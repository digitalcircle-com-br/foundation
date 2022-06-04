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



func Req(c context.Context) *http.Request {
	raw := c.Value(model.CTX_REQ)
	return raw.(*http.Request)
}

func Res(c context.Context) http.ResponseWriter {
	raw := c.Value(model.CTX_RES)
	return raw.(http.ResponseWriter)
}

func SessionID(c context.Context) string {
	ck, err := Req(c).Cookie(string(model.COOKIE_SESSION))
	if err != nil {
		return ""
	}
	return ck.Value
}

func Tenant(c context.Context) string {
	sess := Session(c)
	if sess == nil {
		return ""
	}

	return sess.Tenant
}

func Db(c context.Context) (ret *gorm.DB, err error) {

	t := Tenant(c)
	if t == "" {
		return nil, errors.New("tenant not found")
	}
	ret, err = dbmgr.DBN(t)
	return
}

func Vars(c context.Context) map[string]string {
	return mux.Vars(Req(c))
}

func Done(c context.Context) func() {
	raw := c.Value(model.CTX_DONE)
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
