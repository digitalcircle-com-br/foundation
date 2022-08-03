package sessionmgr

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/migration"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// se sessionKey will generate one string with tenant and session id.
func sessionKey(t string, id string) string {
	return fmt.Sprintf("session:%s:%s", t, id)
}

// sessionKeyFromId will parse session from rawid string.
func sessionKeyFromId(rawid string) (t string, sid string, hash []byte, err error) {
	rawdec, err := base64.StdEncoding.DecodeString(rawid)
	if err != nil {
		return
	}

	parts := strings.Split(string(rawdec), ".")
	if len(parts) != 3 {
		err = errors.New("session id in wrong format")
		return
	}
	t = parts[0]
	sid = parts[1]

	hash, err = base64.StdEncoding.DecodeString(parts[2])
	return
}

// SessionSave - persists the session.
func SessionSave(s *model.Session) (id string, err error) {
	sid := uuid.NewString()
	s.Sessionid = sid
	sessbs, _ := json.Marshal(s)
	hash := md5.New()
	hash.Write(sessbs)
	sum := hash.Sum(nil)
	hashEnc := base64.StdEncoding.EncodeToString(sum)
	iddec := fmt.Sprintf("%s.%s.%s", s.Tenant, s.Sessionid, hashEnc)
	id = base64.StdEncoding.EncodeToString([]byte(iddec))

	k := sessionKey(s.Tenant, s.Sessionid)

	rawSess := model.RawSession{Id: k, Data: sessbs}
	err = db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rawSess).Error

	return
}

func SessionLoad(rawid string) (sess *model.Session, err error) {
	var rawSession model.RawSession
	t, id, hash, err := sessionKeyFromId(rawid)
	if err != nil {
		return
	}

	k := sessionKey(t, id)
	err = db.Where("id = ?", k).First(&rawSession).Error

	if err != nil {
		return nil, err
	}
	hasher := md5.New()
	hasher.Write(rawSession.Data)
	hashVal := hasher.Sum(nil)

	if !bytes.Equal(hash, hashVal) {
		err = errors.New("session hash invalid")
		return
	}

	ret := &model.Session{}
	err = json.Unmarshal(rawSession.Data, ret)
	return ret, err
}

func SessionDel(rawid string) (err error) {
	t, id, _, err := sessionKeyFromId(rawid)
	if err != nil {
		return
	}
	k := sessionKey(t, id)
	return db.Where("id = ?", k).Delete(&model.RawSession{}).Error

}
func SessionDelTenantAndId(t, id string) (err error) {
	k := sessionKey(t, id)
	return db.Where("id = ?", k).Delete(&model.RawSession{}).Error
}

func SessionEnc(s *model.Session) (id string, sessbs []byte) {
	sessbs, _ = json.Marshal(s)
	hasher := md5.New()
	hasher.Write(sessbs)
	sum := hasher.Sum(nil)
	hashEnc := base64.StdEncoding.EncodeToString(sum)
	iddec := fmt.Sprintf("%s.%s.%s", s.Tenant, s.Sessionid, hashEnc)
	id = base64.StdEncoding.EncodeToString([]byte(iddec))
	return
}

func SessionDec(s string) (t string, id string, hash []byte, err error) {
	rawdec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}
	parts := strings.Split(string(rawdec), ".")
	if len(parts) != 3 {
		err = errors.New("session id in wrong format")
		return
	}
	t = parts[0]
	id = parts[1]

	hash, err = base64.StdEncoding.DecodeString(parts[2])
	return

}

var db *gorm.DB

func Setup(d *gorm.DB) error {
	db = d
	return migration.Run(db, migration.Mig{Id: "session-001", Up: func(db *gorm.DB) error {
		return db.AutoMigrate(model.RawSession{})
	},
	})
}
