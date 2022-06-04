package f8n

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"github.com/brianvoe/gofakeit/v6"
// 	"github.com/digitalcircle-com-br/buildinfo"
// 	"github.com/gorilla/mux"
// )

// var server *http.Server
// var router *mux.Router

// func HttpRun(addr string) {
// 	HttpStart(addr)
// 	LockMainRoutine()
// }

// func HttpStart(addr string) *http.Server {
// 	if addr == "" {
// 		addr = ":8080"
// 	}
// 	Log("Server will listen at: %s", addr)
// 	server = &http.Server{Addr: addr, Handler: router}
// 	go func() {
// 		err := server.ListenAndServe()
// 		if err != nil {
// 			Err(err)
// 		}
// 		server = nil
// 	}()
// 	return server
// }

// func HttpStop() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	server.Shutdown(ctx)
// 	server = nil
// }

// func HttpRouter() *mux.Router {
// 	if router == nil {
// 		router = mux.NewRouter()

// 		router.Path("/__test").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte(buildinfo.String()))
// 		})

// 		router.Path("/__help").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Add("Content-type", "application/json")
// 			json.NewEncoder(w).Encode(apiEntries)
// 		})
// 	}
// 	return router
// }

// var apiEntries = make([]*ApiEntry, 0)

// func HttpHandle[TIN, TOUT any](hpath string, method string, perm PermDef, f func(context.Context, TIN) (TOUT, error)) {
// 	Log("Adding handler: %s:%s[%s]", method, hpath, perm)

// 	entry := &ApiEntry{
// 		Path:   hpath,
// 		Method: method,
// 		Perm:   perm,
// 		In:     new(TIN),
// 		Out:    new(TOUT),
// 	}

// 	gofakeit.Struct(entry.In)
// 	gofakeit.Struct(entry.Out)

// 	apiEntries = append(apiEntries, entry)

// 	router.Path(hpath).Methods(method).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		done := false
// 		nctx := context.WithValue(r.Context(), CTX_REQ, r)
// 		nctx = context.WithValue(nctx, CTX_RES, w)
// 		nctx = context.WithValue(nctx, CTX_DONE, func() {
// 			done = true
// 		})

// 		if perm != PERM_ALL {
// 			ck, err := r.Cookie(string(COOKIE_SESSION))
// 			if err != nil {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			sess, err := SessionLoad(ck.Value)
// 			if err != nil || sess == nil {
// 				http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 				return
// 			}
// 			if perm != PERM_AUTH {
// 				_, ok := sess.Perms[perm]
// 				if !ok {
// 					_, ok = sess.Perms[PERM_ROOT]
// 					if !ok {
// 						http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 						return
// 					}
// 				}
// 			}
// 			nctx = context.WithValue(nctx, CTX_SESSION, sess)

// 		}

// 		r = r.WithContext(nctx)
// 		in := new(TIN)

// 		switch r.Method {
// 		case http.MethodPatch:
// 			fallthrough
// 		case http.MethodPost:
// 			fallthrough
// 		case http.MethodPut:
// 			err := json.NewDecoder(r.Body).Decode(in)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 		default:

// 		}

// 		out, err := f(r.Context(), *in)
// 		if !done {
// 			w.Header().Add("Content-Type", "application/json")
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			json.NewEncoder(w).Encode(out)
// 		}

// 	})
// }
