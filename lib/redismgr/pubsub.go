package redismgr

import (
	"context"

	"github.com/digitalcircle-com-br/foundation/lib/core"
)

type Msg struct {
	Chan    string
	Payload string
	Err     error
}

func RawSub(ch ...string) (ret <-chan string, closefn func()) {
	sub := Cli().Subscribe(context.Background(), ch...)
	inret := make(chan string)
	ret = inret
	run := true
	go func() {
		for run {
			msg, err := sub.ReceiveMessage(context.Background())
			if !run {
				return
			}
			if err != nil {
				core.Err(err)
				continue
			}
			inret <- msg.Payload
		}
	}()

	closefn = func() {
		run = false
		sub.Close()
		close(inret)
	}
	return
}

func Sub(ch ...string) (ret <-chan *Msg, closefn func()) {
	sub := Cli().Subscribe(context.Background(), ch...)
	inret := make(chan *Msg)
	ret = inret
	run := true
	go func() {
		for run {
			msg, err := sub.ReceiveMessage(context.Background())
			m := &Msg{}
			if err == nil {
				m.Chan = msg.Channel
				m.Payload = msg.Payload
			} else {
				m.Err = err
			}
			if run {
				inret <- m
			}
		}
	}()

	closefn = func() {
		run = false
		sub.Close()
		close(inret)
	}
	return
}

func Pub(ch string, msg interface{}) {
	Cli().Publish(context.Background(), ch, msg)
}
