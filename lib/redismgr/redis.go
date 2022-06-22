package redismgr

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	libredis "github.com/go-redis/redis/v8"
)

var rediscli *libredis.Client

func Cli() *libredis.Client {
	if rediscli == nil {
		redisurl := os.Getenv("REDIS")
		if redisurl == "" {
			redisurl = "redis://redis:6379"
		}
		opts, err := libredis.ParseURL(redisurl)
		if err != nil {
			panic(err)
		}
		i := 1
		for {
			rediscli = libredis.NewClient(opts)

			context, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, err = rediscli.Ping(context).Result()

			if err == nil {
				return rediscli
			} else {
				core.Warn("could not connect to redis - will retry (%v/10)", i)
				time.Sleep(time.Second)
				i++
				if i >= 10 {
					panic(err)
				}
			}
		}
	}
	return rediscli
}

func HGet(k string, v string) (string, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().HGet(ctx, k, v)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}

func HGetAll(k string) (map[string]string, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().HGetAll(ctx, k)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Result()
}

func Get(k string, i ...interface{}) (string, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Get(ctx, fmt.Sprintf(k, i...))
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}
func GetI(k string, i ...interface{}) (int64, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Get(ctx, fmt.Sprintf(k, i...))
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	str, err := cmd.Result()
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseInt(str, 10, 64)
	return ret, err
}

func GetJson(k string, o interface{}, i ...interface{}) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Get(ctx, fmt.Sprintf(k, i...))
	if cmd.Err() != nil && cmd.Err() != libredis.Nil {
		return cmd.Err()
	}
	str, _ := cmd.Result()
	err := json.Unmarshal([]byte(str), o)
	if err != nil {
		return err
	}
	return nil
}

func Set(k string, v string) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Set(ctx, k, v, 0)
	return cmd.Err()
}

func Del(k string) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Del(ctx, k)
	return cmd.Err()
}

func Hset(k string, k2 string, v string) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().HSet(ctx, k, k2, v)
	return cmd.Err()
}

func Keys(p string, i ...interface{}) ([]string, error) {
	k := fmt.Sprintf(p, i...)
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Keys(ctx, k)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Result()
}

func PGet(p string) (map[string]string, error) {
	ks, err := Keys(p)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]string)
	for _, k := range ks {
		v, err := Get(k)
		if err != nil {
			return nil, err
		}
		ret[k] = v
	}
	return ret, nil
}

func Incr(p string) (int64, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	return Cli().Incr(ctx, p).Result()
}

func Decr(p string) (int64, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	return Cli().Decr(ctx, p).Result()
}

func Expire(p string, to time.Duration) (bool, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	return Cli().Expire(ctx, p, to).Result()
}

func ExpireS(p string, to int) (bool, error) {
	return Expire(p, time.Second*time.Duration(to))
}
