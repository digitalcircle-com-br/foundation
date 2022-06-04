package callmgr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/google/uuid"
)

func Do(in *http.Request) (out *http.Response, err error) {
	// rediscli := redismgr.Cli()
	// id := uuid.NewString()

	// in.Header.Set("X-RETURN-QID", id)

	// buf := bytes.Buffer{}
	// in.Write(&buf)

	// context, cancel := core.Ctx()

	// defer cancel()

	// err = rediscli.LPush(context, q, buf.Bytes()).Err()
	// if err != nil {
	// 	return nil, err
	// }

	// cmdret := rediscli.BRPop(context, time.Second*30, id)
	// if cmdret.Err() != nil {
	// 	return nil, cmdret.Err()
	// }
	// strs, err := cmdret.Result()
	// if err != nil {
	// 	return nil, err
	// }

	// out, err = http.ReadResponse(bufio.NewReader(strings.NewReader(strs[1])), in)
	// return
	panic("implement me")
	return nil, nil

}

func DoQ(q string, in *http.Request) (out *http.Response, err error) {
	rediscli := redismgr.Cli()
	id := uuid.NewString()

	in.Header.Set("X-RETURN-QID", id)

	buf := bytes.Buffer{}
	in.Write(&buf)

	context, cancel := core.Ctx()

	defer cancel()

	err = rediscli.LPush(context, q, buf.Bytes()).Err()
	if err != nil {
		return nil, err
	}

	cmdret := rediscli.BRPop(context, time.Second*30, id)
	if cmdret.Err() != nil {
		return nil, cmdret.Err()
	}
	strs, err := cmdret.Result()
	if err != nil {
		return nil, err
	}

	out, err = http.ReadResponse(bufio.NewReader(strings.NewReader(strs[1])), in)
	return

}

func EncQ(q string, in *http.Request) (err error) {
	rediscli := redismgr.Cli()

	buf := bytes.Buffer{}
	in.Write(&buf)

	context, cancel := core.Ctx()

	defer cancel()

	err = rediscli.LPush(context, q, buf.Bytes()).Err()
	if err != nil {
		return err
	}

	return

}

func SimpleEncQ(q string, i interface{}) error {
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "/cmd", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	return EncQ(q, req)
}
