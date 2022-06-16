package crudmgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/dbmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"gorm.io/gorm"
)

const (
	OP_C  = "C"
	OP_R  = "R"
	OP_U  = "U"
	OP_D  = "D"
	OP_AA = "AA"
	OP_AD = "AD"
)

type CrudOpts struct {
	Op           string        `json:"op"`           // OP - can be C R U OR D
	Db           string        `json:"db"`           // REsolved by server, nevermind
	Tb           string        `json:"tb"`           // Table name
	Where        []interface{} `json:"where"`        // Where clause - []interface{}{"id =?" , my_var_id}
	ID           interface{}   `json:"id"`           // When dealing w ID required ops, this is mandatory (update, delete)
	Data         interface{}   `json:"data"`         // Object with data to be managed
	Cols         []interface{} `json:"cols"`         // Define cols returned by select
	Associations []string      `json:"associations"` // See gorm docs - https://gorm.io/docs/
	PageSize     int           `json:"pagesize"`     // Offset for selected records
	Page         int           `json:"page"`

	AssociationTable  string `json:"association_table"`
	AssociationFieldA string `json:"association_field_a"`
	AssociationFieldB string `json:"association_field_b"`
	AssociationIDA    uint   `json:"association_id_a"`
	AssociationIDB    uint   `json:"association_id_b"`
	Debug             bool   `json:"debug"`
	AutoPreload       bool   `json:"auto_preload"`
	//dataObj model.VO
}

type CrudResponse struct {
	Data         interface{} `json:"data"`
	RowsAffected int64       `json:"rowsaffected"`
}

func Retrieve(tp reflect.Type, opts *CrudOpts) (interface{}, error) {
	db, err := dbmgr.DBN(opts.Db)
	if err != nil {
		return nil, err
	}

	if opts.Debug {
		db = db.Debug()
	}

	if opts.PageSize == 0 || opts.PageSize > 1000 {
		opts.PageSize = 1000
	}

	if opts.Page < 1 {
		opts.Page = 1
	}

	offset := (opts.Page - 1) * opts.PageSize

	ret := reflect.MakeSlice(reflect.SliceOf(tp), 0, opts.PageSize).Interface()
	//ret := make([]T, 0)
	tx := db.Table(opts.Tb)

	for _, assoc := range opts.Associations {
		core.Debug("Loading association: %s", assoc)
		tx = tx.Preload(assoc)
	}

	switch {
	case opts.Where != nil && len(opts.Where) == 1:
		tx = tx.Where(opts.Where)
	case opts.Where != nil && len(opts.Where) > 1:
		tx = tx.Where(opts.Where[0], opts.Where[1:]...)
	}

	if opts.AutoPreload {
		tx = tx.Set("gorm:auto_preload", true)
	}

	err = tx.Limit(opts.PageSize).Offset(offset).Find(&ret).Error
	return CrudResponse{Data: ret}, err
}

func Create(opts *CrudOpts) (interface{}, error) {
	db, err := dbmgr.DBN(opts.Db)
	if err != nil {
		return nil, err
	}

	if opts.Debug {
		db = db.Debug()
	}

	err = db.Table(opts.Tb).Create(opts.Data).Error

	return CrudResponse{Data: []interface{}{opts.Data}}, err
}

func Update(opts *CrudOpts) (interface{}, error) {
	db, err := dbmgr.DBN(opts.Db)
	if err != nil {
		return nil, err
	}

	if opts.Debug {
		db = db.Debug()
	}

	tx := db.Table(opts.Tb).Where("id = ?", opts.ID).Updates(opts.Data)
	ret := CrudResponse{
		Data:         nil,
		RowsAffected: tx.RowsAffected,
	}
	return ret, tx.Error
}

func Delete(opts *CrudOpts) (interface{}, error) {
	db, err := dbmgr.DBN(opts.Db)
	if err != nil {
		return nil, err
	}

	if opts.Debug {
		db = db.Debug()
	}

	tx := db.Exec(fmt.Sprintf("delete from %s where id = ?", opts.Tb), opts.ID)

	ret := CrudResponse{
		Data:         nil,
		RowsAffected: tx.RowsAffected,
	}

	return ret, tx.Error
}

func AssociationAssociate(opts *CrudOpts) (interface{}, error) {
	db, err := dbmgr.DBN(opts.Db)
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("insert into \"%s\"(\"%s\",\"%s\") values (?,?)",
		opts.AssociationTable,
		opts.AssociationFieldA,
		opts.AssociationFieldB), opts.AssociationIDA, opts.AssociationIDB).Error
	return nil, err
}
func AssociationDissociate(opts *CrudOpts) (interface{}, error) {
	db, err := dbmgr.DBN(opts.Db)
	if err != nil {
		return nil, err
	}
	err = db.Exec(fmt.Sprintf("delete from into \"%s\" where \"%s\" = ? and \"%s\" = ?",
		opts.AssociationTable,
		opts.AssociationFieldA,
		opts.AssociationFieldB), opts.AssociationIDA, opts.AssociationIDB).Error
	return nil, err
}

func Handle(a interface{}) error {
	tp := reflect.TypeOf(a).Elem()
	db, err := dbmgr.DB()
	if err != nil {
		return err
	}
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(a)
	tb := stmt.Schema.Table

	core.Log("Registering route %s for CRUD %#v", tb, a)

	routemgr.HandleHttp("/crud/"+tb,
		http.MethodPost,
		model.PERM_AUTH,
		func(w http.ResponseWriter, r *http.Request) error {
			sess := ctxmgr.Session(r.Context())
			if sess == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return nil
			}

			opts := new(CrudOpts)

			buf := &bytes.Buffer{}
			io.Copy(buf, r.Body)
			defer r.Body.Close()

			err := json.Unmarshal(buf.Bytes(), opts)
			if err != nil {
				return err
			}

			_, ok := sess.Perms[model.PermDef("crud."+tb+"."+opts.Op)]
			if !ok {
				_, ok = sess.Perms[model.PERM_ROOT]
			}
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return nil
			}

			typeData := reflect.New(tp).Interface()
			bs, _ := json.Marshal(opts.Data)
			json.Unmarshal(bs, typeData)
			opts.Data = typeData
			//no := reflect.New(tp).Interface()

			opts.Db = sess.Tenant
			log.Printf("using tb: %s", tb)
			opts.Tb = tb

			var ret interface{}

			switch opts.Op {

			case OP_C:
				ret, err = Create(opts)
			case OP_R:
				ret, err = Retrieve(tp, opts)
			case OP_U:
				ret, err = Update(opts)
			case OP_D:
				ret, err = Delete(opts)
			case OP_AA:
				ret, err = AssociationAssociate(opts)
			case OP_AD:
				ret, err = AssociationDissociate(opts)
			default:
				http.Error(w, "Unknown op: "+opts.Op, http.StatusBadRequest)
				return nil
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}

			json.NewEncoder(w).Encode(ret)
			return nil
		})

	return nil
}
