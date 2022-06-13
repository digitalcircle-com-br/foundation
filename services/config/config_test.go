package config_test

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/testmgr"
	"github.com/digitalcircle-com-br/foundation/services/config"
)

func TestPost(t *testing.T) {
	w := testmgr.NewInMemResponseWriter()
	r, err := http.NewRequest(http.MethodPost, "/k/k1", strings.NewReader("a1"))
	if err != nil {
		t.Fatal(err)
	}
	config.Setup()
	routemgr.Router().ServeHTTP(w, r)

	if w.Status() != 200 {
		t.Fatalf("Response from POST KEY should be 200, is %v", w.Status())
	}
}

func TestList(t *testing.T) {
	w := testmgr.NewInMemResponseWriter()
	r, err := http.NewRequest(http.MethodGet, "/list", nil)
	if err != nil {
		t.Fatal(err)
	}
	config.Setup()
	routemgr.Router().ServeHTTP(w, r)

	if w.Status() != 200 {
		t.Fatalf("Response from GET list should be 200, is %v", w.Status())
	}
	var files []string
	err = json.NewDecoder(w).Decode(&files)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 1 {
		t.Fatalf("Files len should be at least 1")
	}
}

func TestGet(t *testing.T) {
	w := testmgr.NewInMemResponseWriter()
	r, err := http.NewRequest(http.MethodGet, "/k/k1", nil)
	if err != nil {
		t.Fatal(err)
	}
	config.Setup()
	routemgr.Router().ServeHTTP(w, r)

	if w.Status() != 200 {
		t.Fatalf("Response from GET list should be 200, is %v", w.Status())
	}
	log.Print(string(w.Bytes()))
}

func TestDel(t *testing.T) {
	w := testmgr.NewInMemResponseWriter()
	r, err := http.NewRequest(http.MethodDelete, "/k/k1", nil)
	if err != nil {
		t.Fatal(err)
	}
	config.Setup()
	routemgr.Router().ServeHTTP(w, r)

	if w.Status() != 200 {
		t.Fatalf("Response from GET list should be 200, is %v", w.Status())
	}
	log.Print(string(w.Bytes()))
}
