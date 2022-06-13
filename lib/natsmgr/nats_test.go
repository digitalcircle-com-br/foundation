package natsmgr_test

import (
	"log"
	"testing"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/natsmgr"
	"github.com/nats-io/nats.go"
)

func TestBasic(t *testing.T) {
	natsmgr.Sub("q", func(m *nats.Msg) {
		log.Print(string(m.Data))
		natsmgr.Pub(m.Reply, []byte("OK, Got it"))
	})

	res, err := natsmgr.Req("q", []byte("TEST"), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(string(res))
}

func TestMulti(t *testing.T) {
	natsmgr.Sub("q", func(m *nats.Msg) {
		log.Printf("SUB1: %s", string(m.Data))
		natsmgr.Pub(m.Reply, []byte("OK, Got it1"))
	})
	natsmgr.Sub("q", func(m *nats.Msg) {
		log.Printf("SUB2: %s", string(m.Data))
		natsmgr.Pub(m.Reply, []byte("OK, Got it2"))
	})
	natsmgr.Sub("q", func(m *nats.Msg) {
		log.Printf("SUB3: %s", string(m.Data))
		natsmgr.Pub(m.Reply, []byte("OK, Got it3"))
	})

	res, err := natsmgr.Req("q", []byte("TEST"), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	log.Print(string(res))
}
