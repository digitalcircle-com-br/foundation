package test

// import (
// 	"net/http"

// 	"github.com/digitalcircle-com-br/foundation/lib/core"
// 	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
// 	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
// 	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
// )

// func Run() {
// 	core.Init("test")
// 	routemgr.HandleHttp("/test", http.MethodGet, fmodel.PERM_ALL, func(w http.ResponseWriter, r *http.Request) error {
// 		r.Write(core.LogWriter())
// 		return nil
// 	})

// 	runmgr.RunABlock()
// }
