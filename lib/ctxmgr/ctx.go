package ctxmgr

import (
	"context"
	"errors"
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/dbmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

//Req returns *http.Request stored on ctx under key model.CTX_REQ
func Req(c context.Context) *http.Request {
	raw := c.Value(model.CTX_REQ)
	return raw.(*http.Request)
}

//Res returns http.ResponseWriter stored on ctx under key model.CTX_RES
func Res(c context.Context) http.ResponseWriter {
	raw := c.Value(model.CTX_RES)
	return raw.(http.ResponseWriter)
}

//SessionID returns the session identifier of the request stored on context
func SessionID(c context.Context) string {
	ck, err := Req(c).Cookie(string(model.COOKIE_SESSION))
	if err != nil {
		return ""
	}
	return ck.Value
}

//Tenant returns the tenant defined in model.Session stored on ctx
func Tenant(c context.Context) string {
	sess := Session(c)
	if sess == nil {
		return ""
	}

	return sess.Tenant
}

//Db returns a *gorm.DB based on tenant queried from ctx
func Db(c context.Context) (ret *gorm.DB, err error) {

	t := Tenant(c)
	if t == "" {
		return nil, errors.New("tenant not found")
	}
	ret, err = dbmgr.DBN(t)
	return
}

//Vars returns the mux.Vars from the *http.Request stored on ctx
func Vars(c context.Context) map[string]string {
	return mux.Vars(Req(c))
}

//Done returns a func() stored on ctx under the key model.CTX_DONE
func Done(c context.Context) func() {
	raw := c.Value(model.CTX_DONE)
	return raw.(func())
}

//Err responds the request stored on ctx with err.Error() and status http.StatusInternalServerError
func Err(c context.Context, err error) bool {
	if err != nil {
		core.Err(err)
		http.Error(Res(c), err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

//Session returns *model.Session stored on ctx under the key model.CTX_SESSION
func Session(c context.Context) *model.Session {
	rawsession := c.Value(model.CTX_SESSION)
	if rawsession != nil {
		return rawsession.(*model.Session)
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
