package runmgr

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/gorilla/mux"
)

/*

qserveOnceHttp waits for a request pushed into the list identified by key "queue: ${SvcName}",
processes it with mux.Router, and pushes a with-TTL copy of the response to redis server.
*/
func qserveOnceHttp(ctx context.Context, q string, m *mux.Router) error {
	rediscli := redismgr.Cli()
	cmd := rediscli.BRPop(ctx, time.Second*0, "queue:"+core.SvcName()) // returns a BRPop (Blocking Remove Pop) redis command
	if cmd.Err() != nil {
		return cmd.Err()
	}
	strs, err := cmd.Result() // blocks execution until there is an element to pop and remove from list

	if err != nil {
		return err
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(strs[1]))) // Creates request from element pulled from redis list

	if err != nil {
		return err
	}

	wrt := NewInMemResponseWriter() // Creates response writer

	m.ServeHTTP(wrt, req) // Serves request
	res := http.Response{}
	res.Body = io.NopCloser(wrt.b)
	res.Header = wrt.h
	res.StatusCode = wrt.sc

	buf := bytes.Buffer{}
	err = res.Write(&buf) // Sends empty response to client, is it right ?

	if err != nil {
		return err
	}
	qid := req.Header.Get("X-RETURN-QID")
	if qid != "" {
		err = rediscli.LPush(ctx, "queue:"+qid, buf.Bytes()).Err() // Inserts response into redis list
		rediscli.Expire(ctx, qid, time.Minute)                     // Defines TTL to 1 minute
	}

	return err
}

/*
RunRedis listens for new items in the redis list identified by the key "queue: ${SvcName}".
When a request is pulled from the list, it is processed by the mux.Router and his
response is pushed to redis server.
*/
func RunRedis() error {

	for {
		err := qserveOnceHttp(context.Background(), core.SvcName(), routemgr.Router())
		if err != nil {
			core.Err(err)
			time.Sleep(time.Second)
		}
	}
	return nil
}
