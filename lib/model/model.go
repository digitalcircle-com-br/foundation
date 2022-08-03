package model

import (
	"time"
)

type PK int64

type EMPTY struct{}

type PermDef string

const (
	PERM_ALL     PermDef = ""
	PERM_AUTH    PermDef = "+"
	PERM_ROOT    PermDef = "*"
	PERM_USERADM PermDef = "user-adm"
)

type CtxKey string

const (
	CTX_REQ     CtxKey = "REQ"
	CTX_RES     CtxKey = "RES"
	CTX_DONE    CtxKey = "DONE"
	CTX_ID      CtxKey = "ID"
	CTX_START   CtxKey = "START"
	CTX_SESSION CtxKey = "SESSION"
)

type CookieName string

const (
	COOKIE_SESSION CookieName = "SESSIONID"
	COOKIE_TENANT  CookieName = "TENANTID"
)

type HeaderName string

const (
	HEADER_APIKEY HeaderName = "X-APIKEY"
	HEADER_TENAT  HeaderName = "X-TENANT"
)

type SecUser struct {
	BaseVO
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Hash     string      `json:"-"`
	Tenant   string      `json:"tenant"`
	Enabled  *bool       `json:"enabled"`
	Groups   []*SecGroup `json:"groups" gorm:"many2many:sec_user_groups;"`
}

type SecGroup struct {
	BaseVO
	Name  string     `json:"name"`
	Perms []*SecPerm `json:"perms" gorm:"many2many:sec_group_perms;"`
}

type SecPerm struct {
	BaseVO
	Name string `json:"name"`
	Val  string `json:"val"`
}

type LoginRequest struct {
	Tenant   string `json:"tenant"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
}

type Session struct {
	Username  string             `json:"username"`
	Sessionid string             `json:"sessionid"`
	Perms     map[PermDef]string `json:"perms"`
	Tenant    string             `json:"tenant"`
	CreatedAt time.Time          `json:"created_at"`
	Hash      []byte             `json:"-"`
}

type RawSession struct {
	Id   string
	Data []byte
}
type ApiEntry struct {
	Path   string
	Method string
	Perm   PermDef
	In     interface{}
	Out    interface{}
}

type DBVersion struct {
	ID    string
	Desc  string
	RunAt time.Time
}

func (d DBVersion) TableName() string {
	return "__version"
}

type File struct {
	BaseVO
	Name      string
	Tenant    string
	MimeType  string
	Len       int64
	Owner     string
	Content   []byte
	CreatedAt time.Time
}

type VO interface {
	GetID() uint
}

type BaseVO struct {
	ID PK `json:"id"`
}

// func (b *BaseVO) GetID() uint {
// 	return b.ID
// }
