package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/apiadapter"
	"github.com/digitalcircle-com-br/foundation/lib/authmgr"
	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/digitalcircle-com-br/foundation/lib/sessionmgr"
	"github.com/digitalcircle-com-br/foundation/services/auth/hash"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthOpts struct {
	UseSecure bool
}

var opts AuthOpts

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

func (s *service) Login(ctx context.Context, lr *AuthRequest) (out *fmodel.EMPTY, err error) {

	user := &fmodel.SecUser{}
	if lr == nil {
		return nil, errors.New("request cannot be nil")
	}
	if DB == nil {
		return nil, errors.New("DB cannot be nil")
	}

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
	sess := &fmodel.Session{}
	sess.Tenant = user.Tenant
	sess.Username = user.Username
	sess.Perms = make(map[fmodel.PermDef]string)
	sess.CreatedAt = time.Now()
	for _, gs := range user.Groups {
		for _, p := range gs.Perms {
			sess.Perms[fmodel.PermDef(p.Name)] = p.Val
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
		Path:    "/",
		Domain:  domain,
		Name:    string(fmodel.COOKIE_SESSION),
		Value:   id,
		Expires: time.Now().Add(time.Hour * 24 * 365 * 100),
	}

	if opts.UseSecure {
		ck.Secure = true
		ck.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(res, &ck)
	res.Header().Add("X-TENANT", user.Tenant)
	//ret.SessionID = id
	ret.Tenant = user.Tenant
	return &fmodel.EMPTY{}, nil
}

func (s *service) Logout(ctx context.Context, lr *fmodel.EMPTY) (out bool, err error) {

	req := ctxmgr.Req(ctx)
	res := ctxmgr.Res(ctx)
	domain := strings.Join(strings.Split(req.URL.Hostname(), ".")[1:], ".")
	ck := http.Cookie{
		Path:     "/",
		Domain:   domain,
		Name:     string(fmodel.COOKIE_SESSION),
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

func (s *service) Check(ctx context.Context, lr *fmodel.EMPTY) (out bool, err error) {
	session := ctxmgr.Session(ctx)
	return session != nil, nil
}

func (s *service) CheckPerm(ctx context.Context, lr *fmodel.EMPTY) (out bool, err error) {
	logrus.Infof("Testing this")
	req := ctxmgr.Req(ctx)
	perm := req.URL.Query().Get("perm")
	session := ctxmgr.Session(ctx)
	_, ok := session.Perms[fmodel.PermDef(perm)]
	if !ok {
		_, ok = session.Perms[fmodel.PERM_ROOT]
	}
	return ok, nil
}

func Setup(r *mux.Router, db *gorm.DB, nOpts ...AuthOpts) error {
	if nOpts != nil && len(nOpts) > 0 {
		opts = nOpts[0]
	} else {
		opts = AuthOpts{UseSecure: true}
	}
	if db == nil {
		return errors.New("db cannot be nil")
	}
	var err error
	DB = db
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "auth001",
			Migrate: func(db *gorm.DB) error {

				for _, mod := range []interface{}{
					&fmodel.SecUser{},
					&fmodel.SecGroup{},
					&fmodel.SecPerm{},
				} {
					err := db.AutoMigrate(mod)
					if err != nil {
						return err
					}
				}

				perm := &fmodel.SecPerm{Name: "*", Val: "*"}
				err := db.Create(perm).Error
				if err != nil {
					return err
				}
				group := &fmodel.SecGroup{Name: "root", Perms: []*fmodel.SecPerm{perm}}
				err = db.Create(group).Error
				if err != nil {
					return err
				}
				enabled := true

				user := &fmodel.SecUser{
					Username: "root",
					Hash:     "$argon2id$v=19$m=65536,t=3,p=2$nTPFgXmlMFphn506a/VQ2Q$0Y/KXMMxDb28CzuqGZdShAnNuNs3l3vInJRh3xd5uq4",
					Email:    "root@root.com",
					Tenant:   "foundation",
					Enabled:  &enabled, Groups: []*fmodel.SecGroup{group},
				}

				err = db.Create(user).Error
				if err != nil {
					return err
				}
				return nil
			},
		},
	})
	if m == nil {
		log.Printf("no migration found")
		return nil
	}
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	r.Name("auth.login").Methods(http.MethodPost).Path("/login").Handler(apiadapter.Adapt(Service.Login))
	r.Name("auth.logout").Methods(http.MethodGet).Path("/logout").Handler(apiadapter.Adapt(Service.Logout))
	r.Name("auth.check").Methods(http.MethodGet).Path("/check").Handler(apiadapter.Adapt(Service.Check))
	r.Name("auth.checkperm").Methods(http.MethodGet).Path("/checkperm").Handler(apiadapter.Adapt(Service.CheckPerm))

	authmgr.AddPerm("auth.login", fmodel.PERM_ALL)
	authmgr.AddPerm("auth.logout", fmodel.PERM_AUTH)
	authmgr.AddPerm("auth.check", fmodel.PERM_AUTH)
	authmgr.AddPerm("auth.checkperm", fmodel.PERM_AUTH)

	err = db.Exec("select 1+1").Error
	if err != nil {
		return fmt.Errorf("error acessing db: %s", err.Error())
	}

	return err
}
