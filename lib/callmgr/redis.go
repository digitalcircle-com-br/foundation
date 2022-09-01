package callmgr

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/google/uuid"
)

type RedisCaller struct{}

var rdb *redis.Client

func getCli() (*redis.Client, error) {
	if rdb != nil {
		return rdb, nil
	}
	redisUrl := os.Getenv("REDIS")
	if redisUrl == "" {
		return nil, fmt.Errorf("No REDIS env var set")
	}
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	rdb = redis.NewClient(opts)
	return rdb, nil
}

func (r *RedisCaller) DoQ(q string, in *http.Request) (out *http.Response, err error) {
	rediscli, err := getCli()
	if err != nil {
		return nil, err
	}
	id := uuid.NewString()

	in.Header.Set("X-RETURN-QID", id)

	buf := bytes.Buffer{}
	in.Write(&buf)

	context, cancel := core.Ctx()

	defer cancel()

	err = rediscli.LPush(context, "queue:"+q, buf.Bytes()).Err()
	if err != nil {
		return nil, err
	}

	cmdret := rediscli.BRPop(context, time.Second*30, "queue:"+id)
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

func (r *RedisCaller) EncQ(q string, in *http.Request) (err error) {
	rediscli, err := getCli()
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	in.Write(&buf)

	context, cancel := core.Ctx()

	defer cancel()

	err = rediscli.LPush(context, "queue:"+q, buf.Bytes()).Err()
	if err != nil {
		return err
	}

	return

}
