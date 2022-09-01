package main_test

import (
	"log"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestUp(t *testing.T) {
	nc, err := nats.Connect("nats://localhost:4222") //nats.DefaultURL)
	if err != nil {
		t.Fatal(err)
	}

	nc.Subscribe("some", func(msg *nats.Msg) {
		log.Printf("GOT:%s", string(msg.Data))
		msg.Respond([]byte("OK, got it"))
	})

	rmsg, err := nc.Request("some", []byte("test 123"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Ok, got: %s", string(rmsg.Data))
	log.Printf("Over")
}
