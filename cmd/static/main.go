package main

import (
	"github.com/digitalcircle-com-br/foundation/services/static"
)

func main() {
	err := static.Run()
	if err != nil {
		panic(err)
	}
}
