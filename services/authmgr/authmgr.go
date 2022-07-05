package authmgr

import (
	"context"
	"errors"
	"net/http"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/crudmgr"
	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/dbmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
	"github.com/digitalcircle-com-br/foundation/services/auth/hash"
)

type UpdatePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"ConfirmPassword"`
}

func UpdatePassword(ctx context.Context, request *UpdatePasswordRequest) (interface{}, error) {
	db, err := dbmgr.DBN("auth")

	if err != nil {
		return nil, err
	}

	session := ctxmgr.Session(ctx)

	if session == nil {
		return nil, errors.New("invalid session")
	}

	if request.Password != request.ConfirmPassword {
		return nil, errors.New("password and confirm password must be the same")
	}

	var user model.SecUser

	dbResult := db.Table("sec_users").Where("username = ? AND enabled = true", session.Username).First(&user)

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	oldPasswordIsTheSame, err := hash.Check(request.Password, user.Hash)

	if err != nil {
		return nil, err
	}

	if oldPasswordIsTheSame {
		return nil, errors.New("password cannot be the old password")
	}

	passwordHash, err := hash.Hash(request.Password)

	if err != nil {
		return nil, err
	}

	dbResult = dbResult.Update("hash", passwordHash)

	crudResult := crudmgr.CrudResponse{Data: nil, RowsAffected: dbResult.RowsAffected}

	return crudResult, dbResult.Error
}

func Setup() error {
	crudmgr.SetDefaultTenant("auth")

	crudmgr.MustHandle(&model.SecPerm{})
	crudmgr.MustHandle(&model.SecGroup{})
	crudmgr.MustHandle(&model.SecUser{})

	return nil
}

func Run() error {
	core.Init("authmgr")
	err := Setup()
	if err != nil {
		return err
	}

	routemgr.Handle("/updatepassword", http.MethodPost, model.PERM_AUTH, UpdatePassword)

	routemgr.Router().Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			core.Debug("Got: %s", r.URL.String())
			h.ServeHTTP(w, r)
		})
	})

	err = runmgr.RunABlock()
	return err
}
