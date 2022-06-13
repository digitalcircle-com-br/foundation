package runmgr

import (
	"bufio"
	"bytes"
	"io"
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/natsmgr"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/nats-io/nats.go"
)

func RunNats() error {
	return natsmgr.SubQ(core.SvcName(), func(m *nats.Msg) {

		wrt := NewInMemResponseWriter()

		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(m.Data)))

		if err != nil {
			wrt.sc = 500
			wrt.Write([]byte(err.Error()))
			natsmgr.Pub(m.Reply, wrt.Bytes())
			return
			//return err
		}

		routemgr.Router().ServeHTTP(wrt, req)
		res := http.Response{}
		res.Body = io.NopCloser(wrt.b)
		res.Header = wrt.h
		res.StatusCode = wrt.sc

		buf := bytes.Buffer{}
		res.Write(&buf)

		if err != nil {
			wrt.sc = 500
			wrt.Write([]byte(err.Error()))
			natsmgr.Pub(m.Reply, wrt.Bytes())
			return
		}

		err = natsmgr.Pub(m.Reply, buf.Bytes())
		if err != nil {
			core.Err(err)
		}

	})

}
