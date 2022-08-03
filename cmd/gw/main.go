package main

import (
	"errors"
	"github.com/nats-io/nats.go"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		parts := strings.Split(request.URL.String(), "/")
		svc := parts[1]
		if svc == "api" {
			svc = parts[2]
		}
		var msg nats.Msg = nats.Msg{
			Subject: svc,
			Header:  nats.Header{},
			Data:    nil,
		}
		for k, v := range request.Header {
			for _, vv := range v {
				msg.Header.Add(k, vv)
			}
		}
		msg.Header.Add("HOST", request.Host)
		msg.Header.Add("METHOD", request.Method)
		msg.Header.Add("URL", request.URL.String())
		bs, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		msg.Data = bs
		res, err := nc.RequestMsg(&msg, time.Minute)
		if errors.Is(err, nats.ErrNoResponders) {
			msg.Subject = "default"
			res, err = nc.RequestMsg(&msg, time.Minute)
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		for k, v := range res.Header {
			for _, vv := range v {
				writer.Header().Add(k, vv)
			}
		}
		writer.Write(res.Data)

	})

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
