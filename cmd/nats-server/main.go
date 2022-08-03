package main

import (
	"time"

	stand "github.com/nats-io/nats-streaming-server/server"
)

func main() {

	opts := stand.GetDefaultOptions()

	snopts := stand.NewNATSOptions()

	_, err := stand.RunServerWithOpts(opts, snopts)

	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Minute)
	}
}
