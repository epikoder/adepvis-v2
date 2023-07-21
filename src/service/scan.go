package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/blackjack/webcam"
	"github.com/epikoder/adepvis/src/pkg/fetch"
	"github.com/google/gousb"

	// "github.com/hedhyw/Go-Serial-Detector/pkg/v1/serialdet"
	"github.com/rs/cors"
	"go.bug.st/serial"

	// "github.com/tarm/serial"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	Scanner struct {
		ctx    context.Context
		Chan   chan int
		Killer chan int
		Data   chan string
	}

	Status struct {
		Type    string
		Message string
		Action  string
		Status  bool
	}

	ScanOption struct {
		Camera bool
		Card   bool
	}

	FrameSizes []webcam.FrameSize
)

const (
	SCANNER Channel = "scanner"
)

const (
	V4L2_PIX_FMT_PJPG = 0x47504A50
	V4L2_PIX_FMT_YUYV = 0x56595559
	camAddr           = "/dev/video0"
)

var supportedFormats = map[webcam.PixelFormat]bool{
	V4L2_PIX_FMT_PJPG: true,
	V4L2_PIX_FMT_YUYV: true,
}

func (slice FrameSizes) Len() int {
	return len(slice)
}

// For sorting purposes
func (slice FrameSizes) Less(i, j int) bool {
	ls := slice[i].MaxWidth * slice[i].MaxHeight
	rs := slice[j].MaxWidth * slice[j].MaxHeight
	return ls < rs
}

// For sorting purposes
func (slice FrameSizes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (s *Scanner) Init(ctx context.Context) {
	s.ctx = context.WithValue(ctx, SCANNER, &Status{
		Message: "place card to scan",
		Status:  true,
	})
	runtime.EventsOn(s.ctx, string(SCANNER), func(optionalData ...interface{}) {
		state, ok := optionalData[0].(Status)
		if !ok {
			runtime.LogFatal(s.ctx, "Conversion error")
			return
		}
		if state.Type == "action" && state.Action == "stop" {
			s.Chan <- 1
		}
	})
	s.Chan = make(chan int)
	s.Data = make(chan string)
	s.Killer = make(chan int)
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) StartScanner(option ScanOption) (interface{}, error) {
	runtime.LogDebug(s.ctx, fmt.Sprintf("%v", option))
	if option.Camera {
		return s.scanCamera()
	}
	return s.scanCard()
}

func (s *Scanner) StopScanner() {
	s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
		Type:   "action",
		Action: "stop",
	})
}

func (s *Scanner) SetId(data string) {
	select {
	case s.Data <- data:
		runtime.LogInfo(s.ctx, "ID Configured")
	default:
		runtime.LogError(s.ctx, "is channel blocked?")
	}
}

func (s *Scanner) scanCamera() (interface{}, error) {
	cam, err := webcam.Open(camAddr)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return nil, fmt.Errorf("cannot open camera")
	}
	defer cam.Close()

	var format webcam.PixelFormat = 0
	formatDesc := cam.GetSupportedFormats()

	for f := range formatDesc {
		runtime.LogInfo(s.ctx, fmt.Sprintf("Image formats: %v", f))
	}
	for f := range formatDesc {
		if supportedFormats[f] {
			format = f
			runtime.LogInfo(s.ctx, fmt.Sprintf("Image: %v", format))
			break
		}
	}

	if format == 0 {
		return nil, fmt.Errorf("no avaliable image format")
	}

	frameSizes := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frameSizes)

	var size *webcam.FrameSize = &frameSizes[len(frameSizes)-1]
	if size == nil {
		runtime.LogError(s.ctx, "No matching frame size, exiting")
		return nil, fmt.Errorf("no matching frame size")
	}
	f, w, h, err := cam.SetImageFormat(format, size.MaxWidth, size.MaxHeight)
	if err != nil {
		runtime.LogError(s.ctx, err.Error())
		return nil, err
	}

	if err = cam.StartStreaming(); err != nil {
		// s.Log.Error(err.Error())
		return nil, fmt.Errorf("cannot start camera")
	}
	if err != nil {
		// s.Log.Error(err.Error())
		return nil, fmt.Errorf("cannot create image file")
	}

	var (
		li   chan *bytes.Buffer = make(chan *bytes.Buffer)
		fi   chan []byte        = make(chan []byte)
		back chan struct{}      = make(chan struct{})
	)

	go s.encodeToImage(cam, back, fi, li, w, h, f)

	var srv *http.Server
	ctx := context.Background()
	go func(addr string, li chan *bytes.Buffer) {
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			//remove stale image
			shutdown := false
			srv.RegisterOnShutdown(func() {
				// s.Log.Info("shutdown recieved")
				shutdown = true
			})
			<-li
			const boundary = `frame`
			w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
			multipartWriter := multipart.NewWriter(w)
			multipartWriter.SetBoundary(boundary)
			for {
				img := <-li
				image := img.Bytes()
				iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
					"Content-type":   []string{"image/jpeg"},
					"Content-length": []string{strconv.Itoa(len(image))},
				})
				if err != nil {
					// s.Log.Error(err.Error())
					return
				}
				_, err = iw.Write(image)
				if err != nil {
					// s.Log.Error(err.Error())
					return
				}

				if shutdown {
					return
				}
			}
		})

		var imgFile *os.File
		mux.HandleFunc("/file.jpeg", func(w http.ResponseWriter, r *http.Request) {
			imgFile, _ = os.Create(fmt.Sprintf("%s/adepvis/adepvis-scan.jpeg", os.TempDir()))
			<-li
			img := <-li
			image := img.Bytes()
			_, err := imgFile.Write(image)
			if err != nil {
				runtime.LogError(s.ctx, err.Error())
			}
			http.ServeFile(w, r, imgFile.Name())
		})

		handler := cors.Default().Handler(mux)
		srv = &http.Server{
			Addr:    addr,
			Handler: handler,
		}
		// s.Log.Error(fmt.Sprintf("%v", srv.ListenAndServe()))
	}("localhost:8080", li)

	go func() {
		timeout := uint32(5)
		started := false
		for {
			select {
			case exitCode := <-s.Killer:
				if exitCode == 1 {
					// s.Log.Info("close frame reader")
					return
				}
			default:
				err = cam.WaitForFrame(timeout)
				switch err.(type) {
				case nil:
				case *webcam.Timeout:
					// s.Log.Error(err.Error())
					continue
				default:
					// s.Log.Error(err.Error())
					return
				}

				frame, err := cam.ReadFrame()
				if err != nil {
					// s.Log.Error(err.Error())
					return
				}

				if len(frame) != 0 {
					select {
					case fi <- frame:
						<-back
					default:
					}
				}
				if !started {
					started = true
					s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
						Type:   "action",
						Action: "start",
					})
				}
			}
		}
	}()

	// s.Log.Info("waiting for id / exit")
	shut := func(exit int) {
		if exit == 1 {
			runtime.LogInfo(s.ctx, "shutting down serve")
			if err := srv.Shutdown(ctx); err != nil {
				runtime.LogError(s.ctx, err.Error())
			}
			runtime.LogInfo(s.ctx, "shutdown completed")

			//broadcast
		BROADCAST:
			for n := 0; n <= 2; n++ {
				select {
				case s.Killer <- 1:
					runtime.LogInfo(s.ctx, "sending kill")
				default:
					break BROADCAST
				}
			}
		}
	}
	for {
		select {
		case exitCode := <-s.Chan:
			shut(exitCode)
			return nil, fmt.Errorf("process ended by user")
		case id := <-s.Data:
			shut(1)
			runtime.LogInfo(s.ctx, id)
			return s.fetchData(id)
		}
	}
}

func (s *Scanner) encodeToImage(wc *webcam.Webcam, back chan struct{}, fi chan []byte, li chan *bytes.Buffer, w, h uint32, format webcam.PixelFormat) {
	var (
		frame []byte
		img   image.Image
	)
	defer close(fi)
	defer close(li)
	for {
		select {
		case exitCode := <-s.Killer:
			// s.Log.Info(fmt.Sprintf("code: %d", exitCode))
			if exitCode == 1 {
				// s.Log.Info("close encoder")
				return
			}
		default:
			bframe := <-fi
			// copy frame
			if len(frame) < len(bframe) {
				frame = make([]byte, len(bframe))
			}
			copy(frame, bframe)
			back <- struct{}{}

			switch format {
			case V4L2_PIX_FMT_YUYV:
				yuyv := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
				for i := range yuyv.Cb {
					ii := i * 4
					yuyv.Y[i*2] = frame[ii]
					yuyv.Y[i*2+1] = frame[ii+2]
					yuyv.Cb[i] = frame[ii+1]
					yuyv.Cr[i] = frame[ii+3]

				}
				img = yuyv
			default:
				// s.Log.Error("invalid format ?")
			}
			//convert to jpeg
			buf := &bytes.Buffer{}
			if err := jpeg.Encode(buf, img, nil); err != nil {
				// s.Log.Error(err.Error())
				return
			}
			li <- buf
		}
	}
}

func (s *Scanner) scanCard() (interface{}, error) {
	s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
		Message: "place card on scanner",
		Status:  true,
	})
	var path string
	if list, err := serial.GetPortsList(); err == nil {
		for _, p := range list {
			if strings.Contains(p, "usb-Silicon_Labs_CP2102_USB_to_UART_Bridge_Controller") {
				path = p
			}
		}
	} else {
		err = fmt.Errorf("no device detected")
		s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
			Message: err.Error(),
			Status:  false,
		})
		return nil, err
	}
	ctx := gousb.NewContext()
	defer ctx.Close()

	// config := &serial.Config{Name: path, Baud: 9600, StopBits: 2}
	port, err := serial.Open(path, &serial.Mode{BaudRate: 9600, StopBits: 2})
	if err != nil {
		// s.Log.Error(err.Error())

		err = fmt.Errorf("cannot connect to device")
		s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
			Message: err.Error(),
			Status:  false,
		})
		return nil, err
	}
	s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
		Message: "place card on scanner",
		Status:  true,
	})

	go func() {
		for {
			exitCode := <-s.Chan
			if exitCode == 1 {
				return
			}
		}
	}()
	id, err := readSerial(port, s.ctx)
	if err != nil {
		s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
			Message: func() string {
				if err == io.EOF {
					return "device disconnected"
				}
				return err.Error()
			}(),
			Status: false,
		})
		return nil, err
	}
	s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
		Message: "looking up database...",
		Status:  true,
	})

	if len(fmt.Sprintf("%q", id)) > 13 {
		return nil, fmt.Errorf("remove the card from scanner and try again")
	}
	return s.fetchData(id)
}

func readSerial(p serial.Port, ctx context.Context) (interface{}, error) {
	data := make([]byte, 8)
	defer p.Close()
	timeOut := time.Now().Add(time.Minute * 3).Unix()
	for {
		if time.Now().Unix() > timeOut {
			return nil, fmt.Errorf("device timeout try again")
		}
		buf := make([]byte, 1)
		_, err := p.Read(buf)
		if err != nil {
			runtime.LogError(ctx, err.Error())
			return nil, err
		}

		if strings.Contains(string(buf), "#") {
			break
		}
		data = append(data, buf...)
	}
	return strings.TrimFunc(string(data), func(r rune) bool {
		return !unicode.IsGraphic(r) || !unicode.IsPrint(r)
	}), nil
}

func (s *Scanner) fetchData(id interface{}) (interface{}, error) {
	url := fmt.Sprintf("/officer/find/%s", id)
	// s.Log.Info(url)
	if AuthState == nil {
		// s.Log.Error("session expired")
		s.ctx = context.WithValue(s.ctx, SCANNER, &Status{
			Message: "session expired: login again",
			Status:  false,
		})
		return nil, nil
	}
	res, err := fetch.Get(url, map[string]string{
		"authorization": AuthState.Access,
	})
	if err != nil {
		// s.Log.Error(err.Error())
		return nil, fmt.Errorf("internet error connecting to server")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return string(body), nil
}
