package authmgr

import (
	"context"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var permsMap map[string]fmodel.PermDef = map[string]fmodel.PermDef{}

func AddPerm(s string, p fmodel.PermDef) {
	permsMap[s] = p
}

func GetPerm(rname string) fmodel.PermDef {
	perm, ok := permsMap[rname]
	if !ok {
		if rname == "" {
			return fmodel.PERM_ROOT
		}
		return fmodel.PermDef(rname)
		return perm
	}
	return perm
}

func MWAuthorize(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perm := GetPerm(mux.CurrentRoute(r).GetName())
		logrus.Debugf("Calling: Perm:%s - Route:%s", string(perm), r.URL.String())
		nctx := context.WithValue(r.Context(), fmodel.CTX_REQ, r)
		nctx = context.WithValue(nctx, fmodel.CTX_RES, w)

		if perm != fmodel.PERM_ALL {
			ck, err := r.Cookie(string(fmodel.COOKIE_SESSION))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sess, err := sessionmgr.SessionLoad(ck.Value)
			if err != nil || sess == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if perm != fmodel.PERM_AUTH {
				_, ok := sess.Perms[fmodel.PermDef(perm)]
				if !ok {
					_, ok = sess.Perms[fmodel.PERM_ROOT]
					if !ok {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
				}
			}
			nctx = context.WithValue(nctx, fmodel.CTX_SESSION, sess)

		}
		r = r.WithContext(nctx)

		h.ServeHTTP(w, r)
	})
}
