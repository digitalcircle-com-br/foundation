package routemgr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/gorilla/mux"
)

var router *mux.Router

//Reset unsets *mux.Router
func Reset() {
	router = nil
}

//[Alessandro] -- add CORS support
func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//Router returns mux.Router, if it is nil, it configures one with default handlers
func Router(CORS_ALLOWED_ORIGINS ...*[]string) *mux.Router {
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

				//[Alessandro] -- add CORS support
				fmt.Println("Checking CORS: ", r.Header.Get("Origin"))
				if len(r.Header.Get("Origin")) == 0 {
					fmt.Println("CORS HEADERS NOT WRITTEN (Possibly GET request or same origin, not need headers)...")
				} else if len(CORS_ALLOWED_ORIGINS) == 0 {
					fmt.Println("CORS ALLOWED - ALWAYS ALLOW!!!")
					w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
					w.Header().Add("Access-Control-Allow-Credentials", "true")
				} else if CORS_ALLOWED_ORIGINS[0] == nil {
					fmt.Println("CORS ALLOWED - NIL LIST!!!")
					w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
					w.Header().Add("Access-Control-Allow-Credentials", "true")
				} else if CORS_ALLOWED_ORIGINS[0] != nil && ArrayContains(*CORS_ALLOWED_ORIGINS[0], r.Header.Get("Origin")) {
					fmt.Println("CORS ALLOWED - CUSTOM ALLOWED!!!")
					w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
					w.Header().Add("Access-Control-Allow-Credentials", "true")
				} else {
					fmt.Println("CORS HEADERS NOT WRITTEN (ERROR)")
				}
				h.ServeHTTP(w, r)
			})
		})
	}
	return router
}

func Handle[TIN, TOUT any](hpath string, method string, perm model.PermDef, f func(context.Context, TIN) (TOUT, error),
	CORS_ALLOWED_ORIGINS ...*[]string) {
	var cors *[]string = nil
	if len(CORS_ALLOWED_ORIGINS) > 0 {
		cors = CORS_ALLOWED_ORIGINS[0]
	}
	core.Log("Adding handler: %s:%s[%s]", "QUEUE", hpath, perm)

	Router(cors).Name(string(perm)).Methods(method).PathPrefix(hpath).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

func HandleHttp(hpath string, method string, perm model.PermDef, f func(w http.ResponseWriter, r *http.Request) error,
	CORS_ALLOWED_ORIGINS ...*[]string) {
	var cors *[]string = nil
	if len(CORS_ALLOWED_ORIGINS) > 0 {
		cors = CORS_ALLOWED_ORIGINS[0]
	}
	core.Log("Adding handler: %s:%s[%s]", method, hpath, perm)
	switch {
	case strings.HasSuffix(hpath, "/"):
		Router(cors).Name(string(perm)).Methods(method).PathPrefix(hpath).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := f(w, r)
			if err != nil {
				core.Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	default:
		Router(cors).Name(string(perm)).Methods(method).Path(hpath).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := f(w, r)
			if err != nil {
				core.Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
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
