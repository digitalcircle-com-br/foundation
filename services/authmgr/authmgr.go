package authmgr

import (
	"context"
	"errors"
	"github.com/digitalcircle-com-br/foundation/lib/crudmgr"
	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/fmodel"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/services/auth/hash"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type UpdatePasswordRequest struct {
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

var db *gorm.DB

func UpdatePassword(ctx context.Context, request *UpdatePasswordRequest) (interface{}, error) {

	session := ctxmgr.Session(ctx)

	if session == nil {
		return nil, errors.New("invalid session")
	}

	if request.NewPassword != request.ConfirmPassword {
		return nil, errors.New("password and confirm password must be the same")
	}

	var user fmodel.SecUser

	dbResult := db.Table("sec_users").Where("username = ? AND enabled = true", session.Username).First(&user)

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	oldPasswordIsCorrect, err := hash.Check(request.OldPassword, user.Hash)

	if err != nil {
		return nil, err
	}

	if !oldPasswordIsCorrect {
		return nil, errors.New("invalid password")
	}

	oldPasswordIsTheSame, err := hash.Check(request.NewPassword, user.Hash)

	if err != nil {
		return nil, err
	}

	if oldPasswordIsTheSame {
		return nil, errors.New("password cannot be the old password")
	}

	passwordHash, err := hash.Hash(request.NewPassword)

	if err != nil {
		return nil, err
	}

	dbResult = dbResult.Update("hash", passwordHash)

	crudResult := crudmgr.CrudResponse{Data: nil, RowsAffected: dbResult.RowsAffected}

	return crudResult, dbResult.Error
}

func Setup(r *mux.Router, db *gorm.DB) error {

	crudmgr.MustHandle(r, db, &fmodel.SecPerm{})
	crudmgr.MustHandle(r, db, &fmodel.SecGroup{})
	crudmgr.MustHandle(r, db, &fmodel.SecUser{})

	routemgr.Handle(r, "/changepasswod", http.MethodPost, fmodel.PERM_AUTH, UpdatePassword)
	routemgr.Handle(r, "/changeuserpass", http.MethodPost, fmodel.PERM_USERADM, UpdatePassword)
	routemgr.Handle(r, "/killallsessions", http.MethodPost, fmodel.PERM_USERADM, UpdatePassword)

	return nil
}
