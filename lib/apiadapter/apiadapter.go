package apiadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/authmgr"
	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"

type CorsOpts struct {
	Origins []string
}

func findStrInSlice(s []string, e string) string {
	for _, a := range s {
		if a == e {
			return e
		}
	}
	return ""
}

var opts CorsOpts

func SetCorsOpts(o CorsOpts) {
	opts = o
}
func MWCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if origin := findStrInSlice(opts.Origins, r.Header.Get("Origin")); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
		}

		h.ServeHTTP(w, r)
	})
}

func AdaptErr(h func(response http.ResponseWriter, request *http.Request) error) func(response http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			logrus.Warnf("error processing: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Adapt[TIN, TOUT any](f func(context.Context, TIN) (TOUT, error)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		done := false
		nctx := context.WithValue(r.Context(), fmodel.CTX_DONE, func() {
			done = true
		})

		r = r.WithContext(nctx)
		in := new(TIN)

		switch r.Method {
		case http.MethodPatch:
			fallthrough
		case http.MethodPost:
			fallthrough
		case http.MethodPut:
			err := json.NewDecoder(r.Body).Decode(in)
			if err != nil {
				core.Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:

		}

		out, err := f(r.Context(), *in)
		if !done {
			w.Header().Add("Content-Type", "application/json")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(out)
		}
	})
}

func DumpAPI(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		perm := authmgr.GetPerm(route.GetName())
		methods := "ALL"
		ms, err := route.GetMethods()
		if err == nil {
			if len(ms) > 0 {
				methods = strings.Join(ms, ",")
			}
		}
		ptmpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		rname := route.GetName()
		if rname == "" {
			rname = "N/A"
		}
		logrus.Infof("%s\t=>%s: %s / %s",
			methods,
			ptmpl,
			rname,
			string(perm),
		)
		return nil
	})
}

func Req(c context.Context) *http.Request {
	raw := c.Value(fmodel.CTX_REQ)
	return raw.(*http.Request)
}

func Res(c context.Context) http.ResponseWriter {
	raw := c.Value(fmodel.CTX_RES)
	return raw.(http.ResponseWriter)
}

func Tenant(c context.Context) string {
	sess := Session(c)
	if sess == nil {
		return ""
	}

	return sess.Tenant
}

func Vars(c context.Context) map[string]string {
	return mux.Vars(Req(c))
}

func Done(c context.Context) func() {
	raw := c.Value(fmodel.CTX_DONE)
	return raw.(func())
}

func Err(c context.Context, err error) bool {
	if err != nil {
		core.Err(err)
		http.Error(Res(c), err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func SessionID(c context.Context) string {
	ck, err := Req(c).Cookie(string(fmodel.COOKIE_SESSION))
	if err != nil {
		return ""
	}
	return ck.Value
}

func Session(c context.Context) *fmodel.Session {
	rawsession := c.Value(fmodel.CTX_SESSION)
	if rawsession != nil {
		return rawsession.(*fmodel.Session)
	}
	sid := SessionID(c)
	if sid == "" {
		return nil
	}
	ret, err := sessionmgr.SessionLoad(sid)
	if err != nil {
		return nil
	}
	return ret

}
