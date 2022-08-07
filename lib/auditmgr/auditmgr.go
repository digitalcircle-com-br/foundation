package auditmgr

import (
	"bufio"
	"bytes"
	"net/http"
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/ctxmgr"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuditEntry struct {
	Id        int64
	Tenant    string
	User      string
	CreatedAt time.Time
	Url       string
	Method    string
	Raw       []byte
}

var db *gorm.DB

func Setup(pdb *gorm.DB) error {
	db = pdb
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "audit-0001",
			Migrate: func(db *gorm.DB) error {
				return db.AutoMigrate(AuditEntry{})
			},
		},
	})
	return m.Migrate()
}

func Add(r *http.Request) *http.Request {
	var buf bytes.Buffer
	r.Write(&buf)
	bs := buf.Bytes()

	sess := ctxmgr.Session(r.Context())
	user := "N/A"
	tenant := "N/A"
	if sess != nil {
		user = sess.Username
		tenant = sess.Tenant
	}
	ae := AuditEntry{
		User:      user,
		Tenant:    tenant,
		CreatedAt: time.Now(),
		Url:       r.URL.String(),
		Method:    r.Method,
		Raw:       bs,
	}
	err := db.Create(&ae).Error
	if err != nil {
		logrus.Warnf("Error saving audit log: %s", err.Error())
		return nil
	}
	ret, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(bs)))
	return ret
}

func MWAudit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nr := Add(r)
		if nr == nil {
			http.Error(w, "Error processing audit entry", http.StatusInternalServerError)
			return
		}
		//Since we are leaving session and other stuff here, this is required. d
		nr = nr.WithContext(r.Context())
		h.ServeHTTP(w, nr)
	})
}
