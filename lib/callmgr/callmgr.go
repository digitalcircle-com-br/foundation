package callmgr

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
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
	out, err = caller.DoQ(q, in)
	if err != nil {
		core.Warn("Error enqueuing data at callmgr.DoQ: %s => %v", q, err)
	}
	return

}

func EncQ(q string, in *http.Request) (err error) {
	err = caller.EncQ(q, in)
	if err != nil {
		core.Warn("Error enqueuing data at callmgr.EnqQ: %s => %v", q, err)
	}
	return
}

func SimpleEncQ(q string, i interface{}) error {
	bs, err := json.Marshal(i)
	if err != nil {
		core.Warn("Error marshalling data for enqueueing at callmgr.SimpleEncQ: %s => %v", q, err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "/cmd", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	return EncQ(q, req)
}
