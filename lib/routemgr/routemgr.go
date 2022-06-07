package routemgr

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/gorilla/mux"
)

var router *mux.Router

func Reset() {
	router = nil
}

func Router() *mux.Router {
	if router == nil {
		router = mux.NewRouter()
		router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			core.Log("Not found: %s: %s", r.Method, r.URL.String())
		})
		router.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				perm := model.PermDef(mux.CurrentRoute(r).GetName())
				core.Log("Calling route: %s", r.URL.String())
				nctx := context.WithValue(r.Context(), model.CTX_REQ, r)
				nctx = context.WithValue(nctx, model.CTX_RES, w)

				if perm != model.PERM_ALL {
					ck, err := r.Cookie(string(model.COOKIE_SESSION))
					if err != nil {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}

					sess, err := sessionmgr.SessionLoad(ck.Value)
					if err != nil || sess == nil {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
					if perm != model.PERM_AUTH {
						_, ok := sess.Perms[model.PermDef(perm)]
						if !ok {
							_, ok = sess.Perms[model.PERM_ROOT]
							if !ok {
								http.Error(w, "Unauthorized", http.StatusUnauthorized)
								return
							}
						}
					}
					nctx = context.WithValue(nctx, model.CTX_SESSION, sess)

				}
				r = r.WithContext(nctx)
				h.ServeHTTP(w, r)
			})
		})
	}
	return router
}

func Handle[TIN, TOUT any](hpath string, method string, perm model.PermDef, f func(context.Context, TIN) (TOUT, error)) {
	core.Log("Adding handler: %s:%s[%s]", "QUEUE", hpath, perm)

	Router().Name(string(perm)).Methods(method).PathPrefix(hpath).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		done := false
		nctx := context.WithValue(r.Context(), model.CTX_DONE, func() {
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

func HandleHttp(hpath string, method string, perm model.PermDef, f func(w http.ResponseWriter, r *http.Request) error) {
	core.Log("Adding handler: %s:%s[%s]", method, hpath, perm)
	Router().Name(string(perm)).Methods(method).PathPrefix(hpath).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			core.Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func IfErr(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func SimpleQueueHandle[TIN any](h func(c context.Context, in TIN) error) {
	Handle("/cmd", http.MethodPost, model.PERM_ALL, func(ctx context.Context, in TIN) (out interface{}, err error) {
		err = h(ctx, in)
		if err != nil {
			core.Err(err)
		}
		return
	})
}
