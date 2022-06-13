package natsmgr

import (
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

const Url = "nats://nats:4222"

func Sub(q string, h func(m *nats.Msg)) error {
	var err error
	if nc == nil {
		nc, err = nats.Connect(Url)
		if err != nil {
			core.Warn("error natsmgr.Sub(%s)::Connect:  %s", q, err.Error())
			return err
		}
	}
	_, err = nc.Subscribe(q, h)
	if err != nil {
		core.Warn("error natsmgr.Sub(%s)::Subscribe:%s  %s", q, "queue."+q, err.Error())
		return err
	}
	return err
}

func Pub(q string, bs []byte) error {
	var err error
	if nc == nil {
		nc, err = nats.Connect(Url)
		if err != nil {
			core.Warn("error natsmgr.Pub(%s)::Connect:  %s", q, err.Error())
			return err
		}
	}
	err = nc.Publish(q, bs)
	if err != nil {
		core.Warn("error natsmgr.Pub(%s)::Publish:%s  %s", q, "queue."+q, err.Error())
		return err
	}

	return err
}

func SubQ(q string, h func(m *nats.Msg)) error {
	var err error
	if nc == nil {
		nc, err = nats.Connect(Url)
		if err != nil {
			core.Warn("error natsmgr.SubQ(%s)::Connect:  %s", q, err.Error())
			return err
		}
	}
	_, err = nc.QueueSubscribe(q, "queue."+q, h)
	if err != nil {
		core.Warn("error natsmgr.SubQ(%s)::QueueSubscribe:%s  %s", q, "queue."+q, err.Error())
		return err
	}

	return err
}

func Req(q string, bs []byte, to time.Duration) ([]byte, error) {
	var err error
	if nc == nil {
		nc, err = nats.Connect(nats.DefaultURL)
		if err != nil {
			core.Warn("error natsmgr.Req(%s)::Connect:  %s", q, err.Error())
			return nil, err
		}
	}
	ret, err := nc.Request(q, bs, to)
	if err != nil {
		core.Warn("error natsmgr.Request(%s)::Request: %s", q, "queue."+q, err.Error())
		return nil, err
	}

	return ret.Data, err
}
