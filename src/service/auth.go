package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"time"

	"github.com/epikoder/adepvis/src/pkg/crypt"
	"github.com/epikoder/adepvis/src/pkg/fetch"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	Auth struct {
		ctx context.Context
	}

	User struct {
		First_Name string
		Last_Name  string
	}

	AuthData struct {
		Authenticated bool
		Access        string
		Expires       time.Time
		User          User
	}

	LoginResponse struct {
		Status string
		Data   struct {
			AccessToken  string    `json:"access_token"`
			RefreshToken string    `json:"refresh_token"`
			ExpiresAt    time.Time `json:"expires_at"`
			Role         string    `json:"role"`
			User         User      `json:"user"`
		}
	}

	LoginCredential struct {
		Svn      string `json:"svn"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

const (
	AUTHSTATE Channel = "auth.state"
)

var (
	RememberUser bool      = true
	AuthState    *AuthData = nil
)

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Init(ctx context.Context) {
	a.ctx = context.WithValue(ctx, AUTHSTATE, &AuthData{
		Authenticated: false,
		Access:        "",
		Expires:       time.Now(),
	})

	if RememberUser {
		runtime.EventsOn(a.ctx, string(AUTHSTATE), func(optionalData ...interface{}) {
			// AuthState = state
			path, err := os.UserHomeDir()
			if err != nil {
				runtime.LogFatal(a.ctx, err.Error())
				return
			}

			b, err := json.Marshal(optionalData)
			if err != nil {
				runtime.LogFatal(a.ctx, err.Error())
				return
			}

			s, err := crypt.Base64Encode(string(b), nil)
			if err != nil {
				runtime.LogFatal(a.ctx, err.Error())
				return
			}
			if err := os.Mkdir(path+"/.adepvis", fs.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
				fmt.Println(err)
				runtime.LogFatal(a.ctx, err.Error())
				return
			}
			path = path + "/.adepvis/.private"
			f, err := os.Create(path)
			if err != nil {
				runtime.LogFatal(a.ctx, "file-create:"+err.Error())
				return
			}
			if _, err := f.WriteString(s); err != nil {
				runtime.LogFatal(a.ctx, "file-write:"+err.Error())
			}
		})
	}
}

func (a *Auth) CheckLoginStatus() (*AuthData, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		runtime.LogFatal(a.ctx, err.Error())
		return nil, err
	}

	path = path + "/.adepvis/.private"
	if store, ok := a.ctx.Value(AUTHSTATE).(*AuthData); ok && store.Authenticated {
		runtime.LogInfo(a.ctx, err.Error())
		return store, nil
	} else if RememberUser {
		body, err := ioutil.ReadFile(path)
		if err == nil {
			var data *AuthData
			s, _ := crypt.Base64Decode(string(body), nil)
			if err = json.Unmarshal([]byte(s), &data); err == nil {
				runtime.LogInfo(a.ctx, "Login found")
				b, _ := json.Marshal(data)
				runtime.LogInfo(a.ctx, string(b))
				a.ctx = context.WithValue(a.ctx, AUTHSTATE, data)
				return data, nil
			}
		}
		runtime.LogError(a.ctx, err.Error())
	}

	return nil, fmt.Errorf("unathorized")
}

func (a *Auth) Login(i LoginCredential) (interface{}, error) {
	res, err := fetch.Post("/auth/login_officer", i, nil)
	runtime.LogDebug(a.ctx, res.Status)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case 401:
		return nil, fmt.Errorf("Unauthorized")
	case 400:
		{
			return string(b), nil
		}
	case 200:
		{
			var lR = &LoginResponse{}
			if err := json.Unmarshal(b, lR); err != nil {
				return nil, err
			}

			authState := &AuthData{
				Authenticated: true,
				Access:        "Bearer " + lR.Data.AccessToken,
				Expires:       lR.Data.ExpiresAt,
				User:          lR.Data.User,
			}
			a.ctx = context.WithValue(a.ctx, AUTHSTATE, authState)
			runtime.EventsEmit(a.ctx, string(AUTHSTATE), authState)
			return true, nil
		}
	default:
		return nil, fmt.Errorf("unknown error occured")
	}

}

func (a *Auth) Logout() {
	a.ctx = context.WithValue(a.ctx, AUTHSTATE, &AuthData{
		Authenticated: false,
		Access:        "",
	})
	if RememberUser {
		var (
			path string
			err  error
		)
		if path, err = os.UserHomeDir(); err != nil {
			runtime.LogFatal(a.ctx, err.Error())
			return
		}
		path = path + "/.adepvis/.private"
		if err = os.Remove(path); err != nil {
			runtime.LogFatal(a.ctx, err.Error())
			return
		}
	}
}
