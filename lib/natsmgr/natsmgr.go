package natsmgr

import (
	"time"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

const Url = "nats://nats:4222"

func Sub(q string, h func(m *nats.Msg)) error {
	var err error
	if nc == nil {
		nc, err = nats.Connect(Url)
		if err != nil {
			return err
		}
	}
	_, err = nc.Subscribe(q, h)
	return err
}

func Pub(q string, bs []byte) error {
	var err error
	if nc == nil {
		nc, err = nats.Connect(Url)
		if err != nil {
			return err
		}
	}
	err = nc.Publish(q, bs)
	return err
}

func SubQ(q string, h func(m *nats.Msg)) error {
	var err error
	if nc == nil {
		nc, err = nats.Connect(Url)
		if err != nil {
			return err
		}
	}
	_, err = nc.QueueSubscribe(q, "queue."+q, h)
	return err
}

func Req(q string, bs []byte, to time.Duration) ([]byte, error) {
	var err error
	if nc == nil {
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			return nil, err
		}
	}
	ret, err := nc.Request(q, bs, to)
	return ret.Data, err
}
