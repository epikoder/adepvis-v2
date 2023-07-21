import { Html5Qrcode } from "html5-qrcode";
import { CircularProgress } from "@mui/material";
import { Component } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";

class QrCodeScanner extends Component<
  { onCompleted: (s: string) => void },
  { runing: boolean; completed: boolean }
> {
  state = {
    runing: false,
    completed: false,
  };

  componentDidMount() {
    EventsOn("scanner", async (state) => {
      if (state.Type !== "action") return;
      if (state.Action === "start") {
        this.setState({ runing: true });
        return this.startQueryQR();
      }
    });
  }

  componentDidUpdate() {
    const qrScanner = document.getElementById("qrScanner");
    if (this.state.runing && qrScanner !== null) {
      qrScanner.innerHTML = "";
      const img = document.createElement("img");
      img.src = "";
      img.width = 200;
      img.height = 200;
      qrScanner.append(img);
      setTimeout(() => {
        img.src = "http://localhost:8080";
      }, 500);
    }
  }

  async startQueryQR() {
    let queryRes;
    for (;;) {
      const html5QrCode = new Html5Qrcode("reader");
      const url = "http://localhost:8080/file.jpeg";
      let res;
      try {
        res = await fetch(url);
      } catch (error) {
        console.log(error);
        break;
      }
      if (res.status !== 200) break;
      const imageFile = new File(
        [await res.blob()],
        new Date().getDate().toString(),
        {
          type: "image/jpeg",
        }
      );

      try {
        queryRes = await html5QrCode.scanFileV2(imageFile, false);
        this.setState({ completed: true });
        console.log(queryRes);
        return this.props.onCompleted(queryRes.decodedText);
      } catch (error) {
        console.log(error);
        continue;
      }
    }
  }

  render() {
    return (
      <div>
        {this.state.runing && (
          <div>
            {this.state.completed ? (
              <CircularProgress
                sx={{
                  color: "blue",
                }}
              />
            ) : (
              <div id="qrScanner"></div>
            )}
          </div>
        )}
        <div
          style={{
            visibility: "hidden",
          }}
        >
          <input
            type="file"
            id="reader"
            style={{
              visibility: "hidden",
              width: 0,
            }}
          ></input>
        </div>
      </div>
    );
  }
}

export default QrCodeScanner;
