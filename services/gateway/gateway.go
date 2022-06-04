package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/callmgr"
	"github.com/digitalcircle-com-br/foundation/lib/cfgmgr"
	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type Route struct {
	Src   string `yaml:"src"`
	Dst   string `yaml:"dst"`
	Async bool   `yaml:"async"`
}

func (r Route) String() string {
	return fmt.Sprintf("%s => %s", r.Src, r.Dst)
}

type cfg struct {
	Addr   string  `yaml:"addr"`
	Routes []Route `yaml:"routes"`
}

var Cfg = new(cfg)

var firstLoad = true

func CreateRedirectabledRequest(r *http.Request) (string, error) {
	urlpath := r.URL.Path
	urlpath = strings.Split(urlpath, "?")[0]
	urlparts := strings.Split(urlpath, "/")
	if len(urlparts) < 3 {
		return "static", nil
	}
	q := urlparts[2]
	if q == "" {
		return "static", nil
	}
	parttobereplaced := "/" + strings.Join(urlparts[1:3], "/")
	nurl := strings.Replace(r.URL.String(), parttobereplaced, "", 1)
	nurlo, err := url.Parse(nurl)
	nurlo.Host = r.URL.Host
	if err != nil {
		return "", err
	}
	r.URL = nurlo
	return q, nil

}

func SetupRoute(route Route) {
	var err error
	h := func(w http.ResponseWriter, r *http.Request) {
		originalUrl := r.URL.String()
		r.URL, err = url.Parse(strings.Replace(r.URL.String(), route.Src, "/", 1))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}
		core.Log("Routing: %s to %s:%s", originalUrl, route.Dst, r.URL.String())
		res, err := callmgr.DoQ(route.Dst, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}
		defer res.Body.Close()
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	}

	router.PathPrefix(route.Src).HandlerFunc(h)
}

var router = mux.NewRouter()

func onChange() {
	if Cfg.Addr == "" {
		Cfg.Addr = ":8080"
	}
	if firstLoad {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			routemgr.Router().ServeHTTP(w, r)
		})
		go func() {
			http.ListenAndServe(Cfg.Addr, nil)
		}()
		firstLoad = false
	}

	router = mux.NewRouter()
	for _, route := range Cfg.Routes {

		core.Log("Adding route: %s", route.String())
		SetupRoute(route)
	}

	router.PathPrefix("/api/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalUrl := r.URL.String()
		q, err := CreateRedirectabledRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}

		core.Log("Routing: %s to %s => %s", originalUrl, q, r.URL.String())
		res, err := callmgr.DoQ(q, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}
		defer res.Body.Close()
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	})

	router.PathPrefix("/app/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalUrl := r.URL.String()
		q, err := CreateRedirectabledRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}

		core.Log("Routing: %s to %s => %s", originalUrl, q, r.URL.String())
		res, err := callmgr.DoQ(q, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}
		defer res.Body.Close()
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	})

	router.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		cfgs := make(map[string]interface{})
		err := cfgmgr.Load("client", &cfgs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		k := r.URL.Query().Get("k")
		if k == "" {
			k = "default"
		}
		cfg, ok := cfgs[r.URL.Hostname()]
		if !ok {
			cfg = cfgs[k]
		}
		json.NewEncoder(w).Encode(cfg)
	})

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalUrl := r.URL.String()
		q := "static"

		core.Log("Routing: %s to %s:%s", originalUrl, q, r.URL.String())
		res, err := callmgr.DoQ(q, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			core.Err(err)
			return
		}
		defer res.Body.Close()
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	})
}
func Prepare() error {
	core.Init("gateway")

	go func() {
		err := cfgmgr.Load("routes", Cfg)
		if err != nil && err != redis.Nil {
			panic(err)
		}

		chChange, _, chErr := cfgmgr.UpdateOnChange("routes", Cfg)

		for {
			select {
			case <-chChange:
				onChange()
			case err := <-chErr:
				core.Err(err)
			}

		}
	}()
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		core.Log("No routes set yet for gateway")
		http.NotFound(w, r)
	})
	routemgr.Router().NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})

	return nil
}

func Run() error {
	err := Prepare()
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8080", routemgr.Router())
}
