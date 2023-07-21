package fetch

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

var (
	ApiUrl = "http://localhost:8081"
)

func Post(url string, body interface{}, headers map[string]string) (res *http.Response, err error) {
	var b []byte
	if reflect.TypeOf(body).Kind() != reflect.Array {
		b, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	} else {
		b = body.([]byte)
	}
	req, err := http.NewRequest("POST", ApiUrl+url, strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return
}

func Get(url string, headers map[string]string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", ApiUrl+url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return
}
