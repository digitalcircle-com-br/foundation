package sessionmgr

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/model"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/google/uuid"
)

func sessionKey(t string, id string) string {
	return fmt.Sprintf("session:%s:%s", t, id)
}

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
	err = redismgr.Set(k, string(sessbs))
	return
}

func SessionLoad(rawid string) (sess *model.Session, err error) {

	t, id, hash, err := sessionKeyFromId(rawid)
	if err != nil {
		return
	}

	k := sessionKey(t, id)
	str, err := redismgr.Get(k)

	if err != nil {
		return nil, err
	}
	hasher := md5.New()
	hasher.Write([]byte(str))
	hashVal := hasher.Sum(nil)

	if !bytes.Equal(hash, hashVal) {
		err = errors.New("session hash invalid")
		return
	}

	ret := &model.Session{}
	err = json.Unmarshal([]byte(str), ret)
	return ret, err
}

func SessionDel(rawid string) (err error) {
	t, id, _, err := sessionKeyFromId(rawid)
	if err != nil {
		return
	}
	k := sessionKey(t, id)
	return redismgr.Del(k)
}
func SessionDelTenantAndId(t, id string) (err error) {
	k := sessionKey(t, id)
	return redismgr.Del(k)
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
