package main

import (
	"github.com/digitalcircle-com-br/foundation/services/gateway"
)

func main() {
	err := gateway.Run()
	if err != nil {
		panic(err)
	}
}
