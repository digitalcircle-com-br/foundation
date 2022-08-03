package httpnatsadapter

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
)

func Adapt(nc *nats.Conn, sub string, to time.Duration) func(w http.ResponseWriter, r *http.Request) error {
	if to == 0 {
		to = time.Minute
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		msg := nats.Msg{Subject: sub, Data: bs, Header: nats.Header{}}
		for k, v := range r.Header {
			for _, vv := range v {
				msg.Header.Add(k, vv)
			}
		}
		res, err := nc.RequestMsg(&msg, time.Minute)
		if err != nil {
			return err
		}
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		sc := http.StatusOK
		//OK - this is completely buggy now
		//TODO: Fix this later.
		if res.Header.Get("statuscode") != "" {
			sc, err = strconv.Atoi(res.Header.Get("statuscode"))
			if err != nil {
				sc = http.StatusInternalServerError
			}

		}
		w.WriteHeader(sc)
		w.Write(res.Data)
		return nil
	}

}
