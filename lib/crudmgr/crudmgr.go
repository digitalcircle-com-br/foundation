package crudmgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

// func findAll[T any](db *gorm.DB) ([]T, error) {
// 	var ret = make([]T, 0)
// 	err := db.Find(&ret).Error
// 	return ret, err
// }

// func findByID[T any](db *gorm.DB, id string) (T, error) {
// 	ret := new(T)
// 	err := db.Where("id = ?", id).First(ret).Error
// 	return *ret, err
// }

// func create[T any](db *gorm.DB, t T) error {
// 	return db.Create(t).Error
// }

// func update[T any](db *gorm.DB, t T) error {
// 	return db.Updates(&t).Error
// }

// func delete[T any](db *gorm.DB, id string) error {
// 	return db.Where("id = ?", id).Delete(new(T)).Error
// }

// func Setup[T any](db *gorm.DB, nr *mux.Router, path string, t T) {
// 	logrus.Infof("Setting up route %s", path)

// 	nr.Methods(http.MethodGet).Path(path + "/{id}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		ret, err := findByID[T](db, mux.Vars(request)["id"])
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		json.NewEncoder(writer).Encode(ret)
// 	})
// 	nr.Methods(http.MethodDelete).Path(path + "/{id}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		err := delete[T](db, mux.Vars(request)["id"])
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		json.NewEncoder(writer).Encode("ok")
// 	})
// 	nr.Methods(http.MethodGet).Path(path).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		ret, err := findAll[T](db)
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		json.NewEncoder(writer).Encode(ret)
// 	})
// 	nr.Methods(http.MethodPost).Path(path).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		at := new(T)
// 		err := json.NewDecoder(request.Body).Decode(at)
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		err = create(db, at)
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(writer).Encode(at)
// 	})
// 	nr.Methods(http.MethodPut).Path(path).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		at := new(T)
// 		err := json.NewDecoder(request.Body).Decode(at)
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		err = update(db, at)
// 		if err != nil {
// 			http.Error(writer, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(writer).Encode(at)
// 	})

// }

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

func retrieve[T any](db *gorm.DB, opts *CrudOpts) (interface{}, error) {

	if opts.PageSize == 0 || opts.PageSize > 1000 {
		opts.PageSize = 1000
	}

	if opts.Page < 1 {
		opts.Page = 1
	}

	offset := (opts.Page - 1) * opts.PageSize

	var ret []T //:= make([]T, 0) //reflect.MakeSlice(reflect.SliceOf(tp), 0, opts.PageSize).Interface()
	//ret := make([]T, 0)
	model := new(T)
	tx := db.Model(model)

	if opts.Debug {
		tx = tx.Debug()
	}

	switch {
	case opts.Where != nil && len(opts.Where) == 1:
		tx = tx.Where(opts.Where)
	case opts.Where != nil && len(opts.Where) > 1:
		tx = tx.Where(opts.Where[0], opts.Where[1:]...)
	}

	switch {

	case opts.AutoPreload:
		tx = tx.Set("gorm:auto_preload", true)

	default:
		for _, assoc := range opts.Associations {
			logrus.Debugf("Loading association: [%s]", assoc)
			tx = tx.Preload(assoc)
		}

	}

	err := tx.Limit(opts.PageSize).Offset(offset).Find(&ret).Error

	return CrudResponse{Data: ret}, err
}

func create(db *gorm.DB, opts *CrudOpts) (interface{}, error) {

	if opts.Debug {
		db = db.Debug()
	}

	err := db.Table(opts.Tb).Create(opts.Data).Error

	return CrudResponse{Data: []interface{}{opts.Data}}, err
}

func update(db *gorm.DB, opts *CrudOpts) (interface{}, error) {

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

func delete(db *gorm.DB, opts *CrudOpts) (interface{}, error) {

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

func associationAssociate(db *gorm.DB, opts *CrudOpts) (interface{}, error) {

	err := db.Exec(fmt.Sprintf("insert into \"%s\"(\"%s\",\"%s\") values (?,?)",
		opts.AssociationTable,
		opts.AssociationFieldA,
		opts.AssociationFieldB), opts.AssociationIDA, opts.AssociationIDB).Error
	return nil, err
}

func associationDissociate(db *gorm.DB, opts *CrudOpts) (interface{}, error) {

	err := db.Exec(fmt.Sprintf("delete from \"%s\" where \"%s\" = ? and \"%s\" = ?",
		opts.AssociationTable,
		opts.AssociationFieldA,
		opts.AssociationFieldB), opts.AssociationIDA, opts.AssociationIDB).Error
	return nil, err
}

func MustHandle[T any](r *mux.Router, db *gorm.DB, a T) {
	err := Handle(r, db, a)
	if err != nil {
		panic(err)
	}
}

var defaultTenant = ""

func SetDefaultTenant(t string) {
	defaultTenant = t
}

func Handle[T any](r *mux.Router, db *gorm.DB, a T) error {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(a)
	tb := stmt.Schema.Table

	logrus.Infof("Registering route %s for CRUD %T", tb, a)

	routemgr.HandleHttp(r, "/"+tb,
		http.MethodPost,
		model.PERM_AUTH,
		func(w http.ResponseWriter, r *http.Request) error {
			defer func() {
				r.Body.Close()
				r := recover()
				if r != nil {
					msg := fmt.Sprintf("Recovering from: %v\n%s ", r, string(debug.Stack()))
					logrus.Warnf(msg)
					http.Error(w, msg, http.StatusInternalServerError)
				}
			}()
			sess := ctxmgr.Session(r.Context())
			if sess == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return nil
			}

			opts := new(CrudOpts)

			buf := &bytes.Buffer{}
			io.Copy(buf, r.Body)

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

			//typeData := reflect.New(tp).Interface()
			typeData := new(T)
			bs, _ := json.Marshal(opts.Data)
			json.Unmarshal(bs, typeData)
			opts.Data = typeData
			//no := reflect.New(tp).Interface()

			if defaultTenant == "" {
				opts.Db = sess.Tenant
			} else {
				opts.Db = defaultTenant
			}
			log.Printf("using tb: %s", tb)
			opts.Tb = tb

			var ret interface{}

			switch opts.Op {

			case OP_C:
				ret, err = create(db, opts)
			case OP_R:
				ret, err = retrieve[T](db, opts)
			case OP_U:
				ret, err = update(db, opts)
			case OP_D:
				ret, err = delete(db, opts)
			case OP_AA:
				ret, err = associationAssociate(db, opts)
			case OP_AD:
				ret, err = associationDissociate(db, opts)
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
