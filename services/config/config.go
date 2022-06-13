package config

import (
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
)

func get(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func post(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func delete(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func list(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func Setup() {
	routemgr.HandleHttp("/list", http.MethodGet, model.PERM_ALL, get)
	routemgr.HandleHttp("/k/", http.MethodGet, model.PERM_ALL, get)
	routemgr.HandleHttp("/k/", http.MethodPost, model.PERM_ALL, post)
	routemgr.HandleHttp("/k/", http.MethodDelete, model.PERM_ALL, delete)
}
