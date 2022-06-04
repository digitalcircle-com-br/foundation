package runmgr

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/gorilla/mux"
)

func Init(s string) {
	core.Init(s)
}

type InMemResponseWriter struct {
	h  http.Header
	b  *bytes.Buffer
	sc int
}

func NewInMemResponseWriter() *InMemResponseWriter {
	ret := &InMemResponseWriter{}
	ret.h = http.Header{}
	ret.b = &bytes.Buffer{}
	ret.sc = http.StatusOK
	return ret
}
func (i *InMemResponseWriter) Header() http.Header {
	return i.h
}
func (i *InMemResponseWriter) Status() int {
	return i.sc
}

func (i *InMemResponseWriter) Write(bs []byte) (int, error) {
	return i.b.Write(bs)
}

func (i *InMemResponseWriter) WriteHeader(statusCode int) {
	i.sc = statusCode
}

func (i *InMemResponseWriter) Read(b []byte) (int, error) {
	return i.b.Read(b)
}
func (i *InMemResponseWriter) Bytes() []byte {
	return i.b.Bytes()
}

func qserveOnceHttp(ctx context.Context, q string, m *mux.Router) error {
	rediscli := redismgr.Cli()
	cmd := rediscli.BRPop(ctx, time.Second*0, core.SvcName())
	if cmd.Err() != nil {
		return cmd.Err()
	}
	strs, err := cmd.Result()

	if err != nil {
		return err
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(strs[1])))

	if err != nil {
		return err
	}

	wrt := NewInMemResponseWriter()

	m.ServeHTTP(wrt, req)
	res := http.Response{}
	res.Body = io.NopCloser(wrt.b)
	res.Header = wrt.h
	res.StatusCode = wrt.sc

	buf := bytes.Buffer{}
	res.Write(&buf)

	if err != nil {
		return err
	}
	qid := req.Header.Get("X-RETURN-QID")
	err = rediscli.LPush(ctx, qid, buf.Bytes()).Err()

	return err
}

func RunA() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	shouldRun := true
	go func() {
		<-ctx.Done()
		shouldRun = false
	}()

	for shouldRun {
		err := qserveOnceHttp(ctx, core.SvcName(), routemgr.Router())
		if err != nil {
			core.Err(err)
			time.Sleep(time.Second)
		}
	}
	return cancel
}

func RunS() context.CancelFunc {
	panic("implement me")
	return nil
}

// func Router() *mux.Router {
// 	return asynchttp.Router()
// }

// func Do(q string, in *http.Request) (out *http.Response, err error) {
// 	return routemgr.Do(q, in)
// }

// func Log(s string, p ...interface{}) {
// 	core.Log(s, p...)
// }

// func Debug(s string, p ...interface{}) {
// 	core.Debug(s, p...)
// }

// func Fatal(s ...interface{}) {
// 	core.Fatal(s...)
// }

// func Err(err error) {
// 	core.Err(err)
// }

// func JsonRead(w http.ResponseWriter, r *http.Request, in interface{}) {
// 	defer r.Body.Close()
// 	err := json.NewDecoder(r.Body).Decode(in)
// 	if err != nil {
// 		Err(err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		in = nil
// 	}
// }

// func JsonWrite(w http.ResponseWriter, r *http.Request, in interface{}) {
// 	err := json.NewEncoder(w).Encode(in)
// 	if err != nil {
// 		Err(err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		in = nil
// 	}
// }

// func DBN(n string) (*gorm.DB, error) {
// 	return dbmgr.DBN(n)
// }

// func HttpIfErr(w http.ResponseWriter, err error) bool {
// 	if err != nil {
// 		Err(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return true
// 	}
// 	return false
// }
