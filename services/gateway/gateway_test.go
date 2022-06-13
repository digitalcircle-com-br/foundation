package gateway_test

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/testmgr"
	"github.com/digitalcircle-com-br/foundation/services/gateway"
	"github.com/stretchr/testify/assert"
)

func TestURLRewrite(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://server.tld:9999/api/some/123", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite2(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/api/some/123", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite3(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/api/some/", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite4(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/api/some", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite5(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/api/", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite6(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/api", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite7(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite8(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestURLRewrite9(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "/asd", nil)
	assert.NoError(t, err)
	q, err := gateway.CreateRedirectabledRequest(r)
	assert.NoError(t, err)
	log.Printf("Q: %s, new URL: %s", q, r.URL.String())
}

func TestConfig(t *testing.T) {
	err := gateway.Prepare()
	time.Sleep(time.Second * 3)
	assert.NoError(t, err)
	r, err := http.NewRequest(http.MethodGet, "/config", nil)
	irw := testmgr.NewInMemResponseWriter()
	routemgr.Router().ServeHTTP(irw, r)
	log.Printf("%#v", r)
	assert.NoError(t, err)
}
func TestConfigKeyA(t *testing.T) {
	err := gateway.Prepare()
	time.Sleep(time.Second * 3)
	assert.NoError(t, err)
	r, err := http.NewRequest(http.MethodGet, "/config?k=a", nil)
	irw := testmgr.NewInMemResponseWriter()
	routemgr.Router().ServeHTTP(irw, r)
	log.Printf("%#v", r)
	assert.NoError(t, err)
}

func TestAppRedir(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://host/app/xpto/a/b/c.txt", nil)
	assert.NoError(t, err)
	_, err = gateway.CreateReverseProxyCall(r, "")
	assert.NoError(t, err)
	log.Printf("new URL: %s", r.URL.String())
}
