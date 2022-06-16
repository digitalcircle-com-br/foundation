package cfgmgr

import (
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"gopkg.in/yaml.v3"
)

func Get(s string) (string, error) {
	return redismgr.Get("config:" + s)
}

func Del(s string) error {
	return redismgr.Del("config:" + s)
}

func Post(s string, body string) error {
	return redismgr.Set("config:"+s, body)
}

func List(s string, body string) ([]string, error) {
	keys, err := redismgr.Keys("config:*")
	if err != nil {
		return keys, err
	}
	var ret []string
	for _, v := range keys {
		ret = append(ret, strings.Replace(v, "config:", "", 1))
	}
	return ret, nil

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

// NotifyChange - informa na mudanca de uma chave de configuracao identificada em s
// retorna um channel que sera disparado qndo da mudanca e uma funcao para fechar o channel.
func NotifyChange(s string) <-chan string {
	ch, _ := redismgr.Sub("config")
	ret := make(chan string)
	go func() {
		for {
			m := <-ch
			ret <- m.Payload
		}
	}()

	return ret
}
