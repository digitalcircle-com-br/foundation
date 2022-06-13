package main

import (
	"github.com/digitalcircle-com-br/foundation/services/config"
)

func main() {
	err := config.Run()
	if err != nil {
		panic(err)
	}
}
