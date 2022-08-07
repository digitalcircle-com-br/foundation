package ctxmgr

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var mx sync.RWMutex

var dbs map[string]*gorm.DB = make(map[string]*gorm.DB)

func AddDb(s string, db *gorm.DB) {
	mx.Lock()
	defer mx.Unlock()
	dbs[s] = db
}

func SetDefaultDB(db *gorm.DB) {
	AddDb("default", db)
}

func DBN(s string) *gorm.DB {
	mx.RLock()
	defer mx.RUnlock()
	return dbs[s]
}

func DB() *gorm.DB {
	return DBN("default")
}

func Req(c context.Context) *http.Request {
	raw := c.Value(fmodel.CTX_REQ)
	return raw.(*http.Request)
}

func Res(c context.Context) http.ResponseWriter {
	raw := c.Value(fmodel.CTX_RES)
	return raw.(http.ResponseWriter)
}

func SessionID(c context.Context) string {
	ck, err := Req(c).Cookie(string(fmodel.COOKIE_SESSION))
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
	ret = DBN(t)
	return
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
