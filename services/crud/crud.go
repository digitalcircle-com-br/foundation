package crud

import (
	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
)

type service struct {
}

var Service = new(service)

func Run() error {
	core.Init("crud")
	runmgr.RunABlock()
	return nil
}
