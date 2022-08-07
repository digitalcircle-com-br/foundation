package apiadapter

import (
	"context"
	"encoding/json"
	"github.com/digitalcircle-com-br/foundation/lib/authmgr"
	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

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
