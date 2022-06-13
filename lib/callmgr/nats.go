package callmgr

import (
	"bufio"
	"bytes"
	"net/http"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/natsmgr"
	"github.com/google/uuid"
)

type NatsCaller struct{}

func (n *NatsCaller) DoQ(q string, in *http.Request) (out *http.Response, err error) {

	id := uuid.NewString()

	in.Header.Set("X-RETURN-QID", id)

	buf := bytes.Buffer{}
	in.Write(&buf)

	ret, err := natsmgr.Req(q, buf.Bytes(), time.Second*30)

	if err != nil {
		return nil, err
	}

	out, err = http.ReadResponse(bufio.NewReader(bytes.NewReader(ret)), in)
	return

}

func (n *NatsCaller) EncQ(q string, in *http.Request) (err error) {

	buf := bytes.Buffer{}
	in.Write(&buf)

	err = natsmgr.Pub(q, buf.Bytes())

	return err

}
