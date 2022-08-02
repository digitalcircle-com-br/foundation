package crud

import (
	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
)

type service struct {
}

var Service = new(service)

/*Run start listening to redis's request queue identified by the key "queue: core.SvcName" */
func Run() error {
	core.Init("crud")
	err := runmgr.RunABlock()
	return err
}
