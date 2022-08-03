package main

import (
	"github.com/digitalcircle-com-br/foundation/services/test"
)

func main() {
	err := test.Run()

	if err != nil {
		panic(err)
	}
}
