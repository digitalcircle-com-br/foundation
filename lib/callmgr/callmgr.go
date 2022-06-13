package callmgr

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Caller interface {
	Do(q string, in *http.Request) (out *http.Response, err error)
	Enc(q string, in *http.Request) (err error)
}

var caller = new(NatsCaller)

var cli http.Client

func Do(in *http.Request) (out *http.Response, err error) {
	return cli.Do(in)
}

func DoQ(q string, in *http.Request) (out *http.Response, err error) {
	return caller.DoQ(q, in)

}

func EncQ(q string, in *http.Request) (err error) {
	return caller.EncQ(q, in)
}

func SimpleEncQ(q string, i interface{}) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "/cmd", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	return EncQ(q, req)
}
