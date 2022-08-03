package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/migration"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/digitalcircle-com-br/foundation/services/auth/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErroNotAuthorized = errors.New("not authorized")

type service struct{}

var Service = new(service)

var DB *gorm.DB

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type AuthResponse struct {
	SessionID string `json:"sessionid"`
	Tenant    string `json:"tenant"`
}

func (s *service) Login(ctx context.Context, lr *AuthRequest) (out *model.EMPTY, err error) {

	user := &model.SecUser{}

	err = DB.Preload("Groups.Perms").Preload(clause.Associations).Where("username = ? and enabled = true", lr.Login).First(user).Error

	if err != nil {
		return nil, err
	}

	match, err := hash.Check(lr.Password, user.Hash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErroNotAuthorized
	}
	sess := &model.Session{}
	sess.Tenant = user.Tenant
	sess.Username = user.Username
	sess.Perms = make(map[model.PermDef]string)
	sess.CreatedAt = time.Now()
	for _, gs := range user.Groups {
		for _, p := range gs.Perms {
			sess.Perms[model.PermDef(p.Name)] = p.Val
		}
	}
	id, err := sessionmgr.SessionSave(sess)
	if err != nil {
		return nil, err
	}
	req := ctxmgr.Req(ctx)
	res := ctxmgr.Res(ctx)
	domain := strings.Join(strings.Split(req.URL.Hostname(), ".")[1:], ".")
	ret := new(AuthResponse)
	ck := http.Cookie{
		Path:     "/",
		Domain:   domain,
		Name:     string(model.COOKIE_SESSION),
		Value:    id,
		Expires:  time.Now().Add(time.Hour * 24 * 365 * 100),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(res, &ck)
	res.Header().Add("X-TENANT", user.Tenant)
	//ret.SessionID = id
	ret.Tenant = user.Tenant
	return &model.EMPTY{}, nil
}

func (s *service) Logout(ctx context.Context, lr *model.EMPTY) (out bool, err error) {

	req := ctxmgr.Req(ctx)
	res := ctxmgr.Res(ctx)
	domain := strings.Join(strings.Split(req.URL.Hostname(), ".")[1:], ".")
	ck := http.Cookie{
		Path:     "/",
		Domain:   domain,
		Name:     string(model.COOKIE_SESSION),
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -24 * 30),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(res, &ck)

	sess := ctxmgr.Session(ctx)
	if sess == nil {
		return true, nil
	}
	err = sessionmgr.SessionDelTenantAndId(sess.Tenant, sess.Sessionid)
	return err == nil, err
}

func (s *service) Check(ctx context.Context, lr *model.EMPTY) (out bool, err error) {
	session := ctxmgr.Session(ctx)
	return session != nil, nil
}

func (s *service) CheckPerm(ctx context.Context, lr *model.EMPTY) (out bool, err error) {
	logrus.Infof("Testing this")
	req := ctxmgr.Req(ctx)
	perm := req.URL.Query().Get("perm")
	session := ctxmgr.Session(ctx)
	_, ok := session.Perms[model.PermDef(perm)]
	if !ok {
		_, ok = session.Perms[model.PERM_ROOT]
	}
	return ok, nil
}

func Setup(r *mux.Router, db *gorm.DB) error {
	var err error
	DB = db
	err = migration.Run(db, migration.Mig{Id: "auth001", Up: func(db *gorm.DB) error {

		for _, mod := range []interface{}{
			&model.SecUser{},
			&model.SecGroup{},
			&model.SecPerm{},
		} {
			err := db.AutoMigrate(mod)
			if err != nil {
				return err
			}
		}

		perm := &model.SecPerm{Name: "*", Val: "*"}
		err := db.Create(perm).Error
		if err != nil {
			return err
		}
		group := &model.SecGroup{Name: "root", Perms: []*model.SecPerm{perm}}
		err = db.Create(group).Error
		if err != nil {
			return err
		}
		enabled := true

		user := &model.SecUser{
			Username: "root",
			Hash:     "$argon2id$v=19$m=65536,t=3,p=2$nTPFgXmlMFphn506a/VQ2Q$0Y/KXMMxDb28CzuqGZdShAnNuNs3l3vInJRh3xd5uq4",
			Email:    "root@root.com",
			Tenant:   "foundation",
			Enabled:  &enabled, Groups: []*model.SecGroup{group},
		}

		err = db.Create(user).Error
		if err != nil {
			return err
		}
		return nil
	},
	})

	routemgr.Handle(r, "/login", http.MethodPost, model.PERM_ALL, Service.Login)
	routemgr.Handle(r, "/logout", http.MethodGet, model.PERM_AUTH, Service.Logout)
	routemgr.Handle(r, "/check", http.MethodGet, model.PERM_AUTH, Service.Check)
	routemgr.Handle(r, "/checkperm", http.MethodGet, model.PERM_AUTH, Service.CheckPerm)

	return err
}

// func Run(db *gorm.DB) error {
// 	var err error
// 	core.Init("auth")
// 	// err = Service.Setup(db)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	routemgr.Handle("/login", http.MethodPost, model.PERM_ALL, Service.Login)
// 	routemgr.Handle("/logout", http.MethodGet, model.PERM_AUTH, Service.Logout)
// 	routemgr.Handle("/check", http.MethodGet, model.PERM_AUTH, Service.Check)
// 	routemgr.Handle("/checkperm", http.MethodGet, model.PERM_AUTH, Service.CheckPerm)

// 	routemgr.Router().Use(func(h http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			logrus.Debugf("Got: %s", r.URL.String())
// 			h.ServeHTTP(w, r)
// 		})
// 	})

// 	err = runmgr.RunABlock()
// 	return err
// }
