package runmgr

import (
	"bytes"
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
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

func RunS() error {
	return http.ListenAndServe(":8080", routemgr.Router())
}

func RunABlock() error {
	return RunRedis()
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
