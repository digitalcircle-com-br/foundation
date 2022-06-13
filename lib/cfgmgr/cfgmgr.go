package cfgmgr

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/natsmgr"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)

func Get(s string) (string, error) {
	res, err := http.Get("http://config:8080/k/" + s)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bs), err
}

func Del(s string) error {
	req, err := http.NewRequest(http.MethodDelete, "http://config:8080/k"+s, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode > 399 {
		return errors.New("response code: " + res.Status)
	}
	return nil

}

func Post(s string, body string) error {
	res, err := http.Post("http://config:8080/k/"+s, "", strings.NewReader(body))
	if err != nil {
		return err
	}

	if res.StatusCode > 399 {
		return errors.New("response code: " + res.Status)
	}
	return nil
}

func List(s string, body string) ([]string, error) {
	res, err := http.Get("http://config:8080/list")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var files []string
	err = json.NewDecoder(res.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	return files, nil
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
func NotifyChange(s string) <-chan string {
	ch := make(chan string)
	natsmgr.Sub("config", func(m *nats.Msg) {
		if string(m.Data) == s {
			str, err := Get(s)
			if err != nil {
				core.Err(err)
				return
			}
			ch <- str
		}
	})
	return ch
}
