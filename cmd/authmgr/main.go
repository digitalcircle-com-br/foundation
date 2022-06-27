package main

import (
	"github.com/digitalcircle-com-br/foundation/services/authmgr"
)

func main() {
	err := authmgr.Run()
	if err != nil {
		panic(err)
	}
}
