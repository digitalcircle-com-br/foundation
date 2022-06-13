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

func qserveOnceHttp(ctx context.Context, q string, m *mux.Router) error {
	rediscli := redismgr.Cli()
	cmd := rediscli.BRPop(ctx, time.Second*0, core.SvcName())
	if cmd.Err() != nil {
		return cmd.Err()
	}
	strs, err := cmd.Result()

	if err != nil {
		return err
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(strs[1])))

	if err != nil {
		return err
	}

	wrt := NewInMemResponseWriter()

	m.ServeHTTP(wrt, req)
	res := http.Response{}
	res.Body = io.NopCloser(wrt.b)
	res.Header = wrt.h
	res.StatusCode = wrt.sc

	buf := bytes.Buffer{}
	res.Write(&buf)

	if err != nil {
		return err
	}
	qid := req.Header.Get("X-RETURN-QID")
	if qid != "" {
		err = rediscli.LPush(ctx, qid, buf.Bytes()).Err()
		rediscli.Expire(ctx, qid, time.Minute)
	}

	return err
}

func RunRedis() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	shouldRun := true
	go func() {
		<-ctx.Done()
		shouldRun = false
	}()

	for shouldRun {
		err := qserveOnceHttp(ctx, core.SvcName(), routemgr.Router())
		if err != nil {
			core.Err(err)
			time.Sleep(time.Second)
		}
	}
	return cancel
}
