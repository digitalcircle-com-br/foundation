package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/digitalcircle-com-br/foundation/lib/dbmgr"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/routemgr"
	"github.com/digitalcircle-com-br/foundation/lib/runmgr"
	"gorm.io/gorm"
)

type service struct{}

var Service = new(service)

func (s service) Setup() error {
	dbmgr.MigrationAdd("files-0001", "Creates Basic Files Infra", func(s string) bool {
		return s != "postgres" && s != "auth"
	}, func(d *gorm.DB) error {
		return d.AutoMigrate(&model.File{})
	})

	return dbmgr.MigrationRun()
}

type UploadResponseEntry struct {
	Id        uint
	Filename  string
	Fieldname string
}

func (s service) Upload(w http.ResponseWriter, r *http.Request) {
	sess := ctxmgr.Session(r.Context())
	if sess == nil {
		return
	}
	db, err := ctxmgr.Db(r.Context())
	if err != nil {
		return
	}
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	mp := r.MultipartForm
	if mp == nil || mp.File == nil || len(mp.File) < 1 {
		return
	}
	resp := make([]UploadResponseEntry, 0)
	for fieldName, v := range r.MultipartForm.File {
		for _, vv := range v {
			mt := mime.TypeByExtension(filepath.Ext(vv.Filename))
			uploadedFile, err := vv.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			buf := &bytes.Buffer{}
			defer uploadedFile.Close()
			io.Copy(buf, uploadedFile)

			f := &model.File{
				Name:      vv.Filename,
				Len:       vv.Size,
				Owner:     sess.Username,
				Tenant:    sess.Tenant,
				MimeType:  mt,
				Content:   buf.Bytes(),
				CreatedAt: time.Now()}

			err = db.Save(f).Error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resp = append(resp, UploadResponseEntry{Id: uint(f.ID), Filename: vv.Filename, Fieldname: fieldName})
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func (s service) Download(w http.ResponseWriter, r *http.Request) {
	sess := ctxmgr.Session(r.Context())

	if sess == nil {
		return
	}

	db, err := ctxmgr.Db(r.Context())
	if err != nil {
		return
	}

	id := r.URL.Query().Get("f")
	attachement := r.URL.Query().Get("attachement")
	if id == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	f := &model.File{}
	err = db.Where("id = ?", id).First(f).Error

	if err == gorm.ErrRecordNotFound {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Description", "File Transfer")
	if attachement == "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", f.Name))
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", f.Name))
	}
	w.Header().Set("Content-Type", f.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%v", f.Len))
	w.WriteHeader(http.StatusOK)
	w.Write(f.Content)
}

func (s service) List(w http.ResponseWriter, r *http.Request) {
	files := make([]model.File, 0)

	sess := ctxmgr.Session(r.Context())

	if sess == nil {
		return
	}

	db, err := ctxmgr.Db(r.Context())
	if routemgr.IfErr(w, err) {
		return
	}

	tx := db.Select("id", "name", "mime_type", "created_at", "owner").Where("tenant = ?", sess.Tenant)

	name := r.URL.Query().Get("name")
	if name != "" {
		tx = tx.Where("name like ?", name)
	}

	dtini := r.URL.Query().Get("dtini")
	dtend := r.URL.Query().Get("dtend")

	switch {
	case dtini != "" && dtend != "":
		tx = tx.Where("created_at between ? and ?", dtini, dtend)
	case dtini != "":
		tx = tx.Where("created_at > ?", dtini)
	case dtend != "":
		tx = tx.Where("created_at < ?", dtend)
	}

	err = tx.Find(&files).Error
	if routemgr.IfErr(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)

}
func Run() error {
	core.Init("files")
	err := Service.Setup()
	if err != nil {
		return err
	}
	routemgr.Router().Name("file.upload").Methods(http.MethodPost).Path("upload").HandlerFunc(Service.Download)
	routemgr.Router().Name("file.download").Methods(http.MethodGet).Path("download").HandlerFunc(Service.Download)
	routemgr.Router().Name("file.list").Methods(http.MethodPost).Path("list").HandlerFunc(Service.List)
	routemgr.Router().Name("file.del").Methods(http.MethodPost).Path("del").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	runmgr.RunABlock()

	return nil

}
