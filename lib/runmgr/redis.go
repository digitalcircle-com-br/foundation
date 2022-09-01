package runmgr

// func qserveOnceHttp(ctx context.Context, q string, m *mux.Router) error {
// 	rediscli := redismgr.Cli()
// 	cmd := rediscli.BRPop(ctx, time.Second*0, "queue:"+core.SvcName())
// 	if cmd.Err() != nil {
// 		return cmd.Err()
// 	}
// 	strs, err := cmd.Result()

// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(strs[1])))

// 	if err != nil {
// 		return err
// 	}

// 	wrt := NewInMemResponseWriter()

// 	m.ServeHTTP(wrt, req)
// 	res := http.Response{}
// 	res.Body = io.NopCloser(wrt.b)
// 	res.Header = wrt.h
// 	res.StatusCode = wrt.sc

// 	buf := bytes.Buffer{}
// 	res.Write(&buf)

// 	if err != nil {
// 		return err
// 	}
// 	qid := req.Header.Get("X-RETURN-QID")
// 	if qid != "" {
// 		err = rediscli.LPush(ctx, "queue:"+qid, buf.Bytes()).Err()
// 		rediscli.Expire(ctx, qid, time.Minute)
// 	}

// 	return err
// }

// func RunRedis() error {
// 	for {
// 		err := qserveOnceHttp(context.Background(), core.SvcName(), routemgr.Router())
// 		if err != nil {
// 			core.Err(err)
// 			time.Sleep(time.Second)
// 		}
// 	}
// 	return nil
// }
