package f8n

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"path/filepath"
// 	"syscall"
// 	"time"

// 	"github.com/digitalcircle-com-br/buildinfo"
// 	"github.com/go-redis/redis/v8"
// 	"gopkg.in/yaml.v2"
// )

// type EMPTY_TYPE struct{}

// func Ctx() (context.Context, context.CancelFunc) {
// 	return context.WithTimeout(context.Background(), time.Second*120)
// }

// func IsDocker() bool {
// 	_, err := os.Stat("/.dockerenv")
// 	return err == nil
// }

// var svcName = ""
// var sigCh = make(chan os.Signal)
// var rediscli *redis.Client

// var onStop = func() {
// 	Log("Terminating")
// }

// func RedisCli() *redis.Client {
// 	return rediscli
// }

// func Init(s string) {
// 	svcName = s
// 	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

// 	redisurl := os.Getenv("REDIS")
// 	if redisurl == "" {
// 		redisurl = "redis://redis:6379"
// 	}
// 	opts, err := redis.ParseURL(redisurl)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rediscli = redis.NewClient(opts)

// 	context, cancel := Ctx()

// 	defer cancel()

// 	_, err = rediscli.Ping(context).Result()

// 	if err != nil {
// 		//TODO: improve error msg
// 		panic(err)
// 	}

// 	go func() {
// 		<-sigCh
// 		err := rediscli.Close()
// 		Err(err)
// 		onStop()
// 		if server != nil {
// 			HttpStop()
// 		}
// 		os.Exit(0)
// 	}()
// 	if IsDocker() {
// 		Log("Initiating Container for: %s", svcName)
// 	} else {
// 		Log("Initiating Service: %s", svcName)
// 	}
// 	wd, _ := os.Getwd()
// 	abswd, _ := filepath.Abs(wd)
// 	Log("Running from %s", abswd)
// 	Log(buildinfo.String())
// }

// var cfg = []byte{}

// func Config(i interface{}) chan struct{} {
// 	ret := make(chan struct{})
// 	go func() {
// 		for {
// 			lastCfg := cfg
// 			cfgstr, err := DataGet(svcName)

// 			if err == nil {
// 				cfgbs := []byte(cfgstr)
// 				if !bytes.Equal(cfgbs, lastCfg) {
// 					cfg = cfgbs
// 					yaml.Unmarshal(cfg, i)
// 					ret <- struct{}{}
// 				}
// 			}

// 			time.Sleep(time.Duration(10) * time.Second)

// 		}
// 	}()
// 	return ret
// }

// func OnStop(f func()) {
// 	onStop = f
// }

// func LockMainRoutine() {
// 	for {
// 		time.Sleep(time.Minute)
// 	}
// }

// func ServerTiming(w http.ResponseWriter, metric string, desc string, t time.Time) {
// 	dur := time.Since(t)
// 	v := float64(dur.Nanoseconds()) / float64(1000000)
// 	w.Header().Add("Server-Timing", fmt.Sprintf("%s;desc=\"%s\";dur=%v", metric, desc, v))
// 	Debug("Server time: %s(%s): %v", desc, metric, v)
// }
