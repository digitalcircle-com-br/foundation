package main

import (
	"errors"
)

var cmd struct {
	Load     struct{} `cmd:""`
	Run      struct{} `cmd:""`
	Clean    struct{} `cmd:""`
	Drop     struct{} `cmd:""`
	All      struct{} `cmd:""`
	Hashpass struct {
		In string `help:"Pass to Hash."`
	} `cmd:"hashpass"`
}

func main() {
	panic(errors.New("implement me"))
	// ctx := kong.Parse(&cmd)
	// var err error

	// switch ctx.Command() {
	// case "load":
	// 	err = setup.Load()
	// case "run":
	// 	err = setup.Run()
	// case "clean":
	// 	err = setup.Clean()
	// case "drop":
	// 	err = setup.Drop()
	// case "hashpass":
	// 	out, err := setup.CreatePassHash(cmd.Hashpass.In)
	// 	if err == nil {
	// 		log.Printf("Hash: %s", out)
	// 	}
	// case "all":
	// 	for _, v := range []func() error{
	// 		setup.Drop,
	// 		setup.Clean,
	// 		setup.Load,
	// 		setup.Run,
	// 	} {
	// 		err = v()
	// 		if err != nil {
	// 			break
	// 		}
	// 	}

	// default:
	// 	err = errors.New("Unknown cmd: " + ctx.Command())
	// }
	// if err != nil {
	// 	panic(err)
	// }

	// log.Printf("Setup done.")
}
