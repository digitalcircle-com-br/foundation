package authmgr

import (
	"github.com/digitalcircle-com-br/foundation/lib/crudmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
)

func Setup() error {

	for _, vo := range []interface{}{
		&model.SecPerm{},
		&model.SecGroup{},
		&model.SecUser{},
	} {
		err := crudmgr.Handle(vo)
		if err != nil {
			return err
		}
	}
	return nil
}
