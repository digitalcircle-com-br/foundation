package testmgr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func Init() {
	core.Init("test")
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

func Get(t *testing.T, k string) string {
	ret, err := redismgr.Get("test:" + k)
	if err == redis.Nil {
		return ""
	}
	if err != nil {
		t.Fatal(fmt.Sprintf("Error getting %s from redis: %s", k, err.Error()))
		return ""
	}
	return ret
}

func Set(t *testing.T, k string, v string) {
	err := redismgr.Set("test:"+k, v)
	if err != nil {
		t.Fatal(fmt.Sprintf("Error setting %s from redis: %s", k, err.Error()))
	}
}

func Login(t *testing.T) {
	// core.Init("auth_test")
	// err := auth.Service.Setup()
	// assert.NoError(t, err)
	// res, err := auth.Service.Login(context.Background(), &auth.AuthRequest{Login: "root", Password: "root"})
	// sessid := res.Cookies()[0].Value
	// assert.NoError(t, err)
	// assert.NotNil(t, res)
	// Set(t, "session", sessid)
	// Set(t, "tenant", res.Tenant)
}

func SessID(t *testing.T) string {
	return Get(t, "session")
}
func Tenant(t *testing.T) string {
	return Get(t, "tenant")
}

func HttpNewAuthRequest(t *testing.T, method string, url string, body []byte, w http.ResponseWriter) *http.Request {
	sessid := SessID(t)
	if sessid == "" {
		Login(t)
		sessid = SessID(t)
		if sessid == "" {
			assert.NoError(t, errors.New("could not setup session"))
		}
	}
	assert.NotEmpty(t, sessid)

	r, err := http.NewRequest(method, url, bytes.NewReader(body))
	assert.NoError(t, err)

	r.Header = http.Header{}
	r.Header.Set("Cookie", fmt.Sprintf("%s=%s", fmodel.COOKIE_SESSION, sessid))
	nctx := context.WithValue(r.Context(), fmodel.CTX_REQ, r)
	nctx = context.WithValue(nctx, fmodel.CTX_RES, w)

	return r
}

func HttpNewAuthRequestO(t *testing.T, method string, url string, body interface{}, w http.ResponseWriter) *http.Request {
	bs, err := json.Marshal(body)
	assert.NoError(t, err)
	return HttpNewAuthRequest(t, method, url, bs, w)
}
