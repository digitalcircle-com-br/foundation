package main

import (
	"github.com/nats-io/nats.go"
	"time"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}

	nc.Subscribe("some", func(msg *nats.Msg) {
		msg.Respond([]byte("Hello World"))
	})

	for {
		time.Sleep(time.Minute)
	}
}
