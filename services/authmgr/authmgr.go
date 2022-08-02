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

	crudmgr.MustHandle(&model.SecPerm{})  // Defines a URI to table sec_perm in mux.Router
	crudmgr.MustHandle(&model.SecGroup{}) // Defines a URI to table sec_group in mux.Router
	crudmgr.MustHandle(&model.SecUser{})  // Defines a URI to table sec_user in mux.Router

	return nil
}

/*Run configures mux.Router and start listening to redis's request queue identified by the key "queue: authmgr" */
func Run() error {
	core.Init("authmgr")
	err := Setup()
	if err != nil {
		return err
	}

	// Middleware to log incoming requests
	routemgr.Router().Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			core.Debug("Got: %s", r.URL.String())
			h.ServeHTTP(w, r)
		})
	})

	err = runmgr.RunABlock() // blocks execution
	return err
}
