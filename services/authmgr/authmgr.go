package authmgr

import (
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/crudmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
)

func Setup() error {
	crudmgr.SetDefaultTenant("auth")

	crudmgr.MustHandle(&model.SecPerm{})
	crudmgr.MustHandle(&model.SecGroup{})
	crudmgr.MustHandle(&model.SecUser{})

	return nil
}

func Run() error {
	core.Init("authmgr")
	err := Setup()
	if err != nil {
		return err
	}
	routemgr.Router().Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			core.Debug("Got: %s", r.URL.String())
			h.ServeHTTP(w, r)
		})
	})

	err = runmgr.RunABlock()
	return err
}
