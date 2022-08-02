package test

import (
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
)

/*Run configures mux.Router and start listening to redis's request queue identified by the key "queue: test" */
func Run() error {
	core.Init("test")
	routemgr.HandleHttp("/test", http.MethodGet, model.PERM_ALL, func(w http.ResponseWriter, r *http.Request) error {
		r.Write(core.LogWriter())
		return nil
	})

	err := runmgr.RunABlock()

	return err
}
