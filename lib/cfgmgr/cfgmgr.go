package cfgmgr

import (
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
)

func Get(s string) (string, error) {
	str, err := redismgr.Get("config:" + s)
	if err == redis.Nil {
		core.Warn("Config %s is not set", s)
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return str, nil
}

func Load(s string, i interface{}) error {
	str, err := Get(s)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(str), i)
	return err
}

func Save(s string, i interface{}) error {
	bs, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	err = redismgr.Set("config:"+s, string(bs))
	if err != nil {
		return err
	}
	redismgr.Pub("config:"+s, string(bs))
	return err
}

//NotifyChange - informa na mudanca de uma chave de configuracao identificada em s
//retorna um channel que sera disparado qndo da mudanca e uma funcao para fechar o channel.
func NotifyChange(s string) (<-chan string, func()) {
	return redismgr.RawSub("config:" + s)
}

var lastcfg = make(map[string]string)

func UpdateOnChange(s string, i interface{}) (chan struct{}, func(), chan error) {
	var cherr = make(chan error)
	var chok = make(chan struct{})
	run := true
	stop := func() {
		run = false
	}
	go func() {
		for run {

			val, err := redismgr.Get("config:%s", s)
			if err == redis.Nil {
				val = ""
				err = nil
			} else if err != nil {
				cherr <- err
				time.Sleep(time.Second * 10)
				continue
			}

			oldv, ok := lastcfg[s]

			if ok && oldv == val {
				time.Sleep(time.Second * 10)
				continue
			}

			lastcfg[s] = val

			if err != nil {
				cherr <- err
				continue
			}

			err = yaml.Unmarshal([]byte(val), i)

			if err != nil {
				cherr <- err
				continue
			} else {
				chok <- struct{}{}
			}
		}
	}()

	return chok, stop, cherr
}
