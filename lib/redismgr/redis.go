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

//Cli returns a redis client
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

			ctx, cancel := context.WithCancel(context.Background())

			_, err = rediscli.Ping(ctx).Result()

			if err == nil {
				cancel()
				return rediscli
			} else {
				cancel()
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

//HGet wraps Cli().HGet() and parse return to string
func HGet(k string, v string) (string, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().HGet(ctx, k, v)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}

//HGetAll wraps Cli().HGetAll() and parse return to string
func HGetAll(k string) (map[string]string, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().HGetAll(ctx, k)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Result()
}

//Get wraps Cli().Get() and parse return to string
func Get(k string, i ...interface{}) (string, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Get(ctx, fmt.Sprintf(k, i...))
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}

//GetI wraps Cli().Get() and parse return to int64
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

//GetJson wraps Cli().Get() and unmarshal return to JSON
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

//Set wraps Cli().Set(), returning its error
func Set(k string, v string) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Set(ctx, k, v, 0)
	return cmd.Err()
}

//Del wraps Cli().Del(), returning its error
func Del(k string) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().Del(ctx, k)
	return cmd.Err()
}

//Hset wraps Cli().Hset(), returning its error
func Hset(k string, k2 string, v string) error {
	ctx, cancel := core.Ctx()
	defer cancel()
	cmd := Cli().HSet(ctx, k, k2, v)
	return cmd.Err()
}

//Keys wraps Cli().Keys()
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

//PGet returns values from all defined keys
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

//Incr wraps Cli().Incr(), returning its result
func Incr(p string) (int64, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	return Cli().Incr(ctx, p).Result()
}

//Decr wraps Cli().Decr(), returning its result
func Decr(p string) (int64, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	return Cli().Decr(ctx, p).Result()
}

//Expire wraps Cli().Expire(), returning its result
func Expire(p string, to time.Duration) (bool, error) {
	ctx, cancel := core.Ctx()
	defer cancel()
	return Cli().Expire(ctx, p, to).Result()
}

//ExpireS calls Expire setting time to "${to} seconds"
func ExpireS(p string, to int) (bool, error) {
	return Expire(p, time.Second*time.Duration(to))
}
