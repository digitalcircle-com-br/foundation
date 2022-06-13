package main

import "github.com/digitalcircle-com-br/foundation/services/auth"

func main() {
	err := auth.Run()
	if err != nil {
		panic(err)
	}
}
