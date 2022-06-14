package dbmgr

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/digitalcircle-com-br/foundation/lib/cfgmgr"
	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbs = make(map[string]*gorm.DB)
var mx sync.RWMutex

func saveDb(n string, d *gorm.DB) {
	mx.Lock()
	defer mx.Unlock()
	dbs[n] = d
}
func loadDb(n string) (d *gorm.DB, ok bool) {
	mx.RLock()
	defer mx.RUnlock()
	d, ok = dbs[n]
	return
}

func delDb(n string) {
	mx.Lock()
	defer mx.Unlock()
	delete(dbs, n)
}

func DB() (ret *gorm.DB, err error) {
	return DBN("foundation")
}
func DBMaster() (ret *gorm.DB, err error) {
	return DBN("postgres")
}

var dsns map[string]string

func DBN(n string) (ret *gorm.DB, err error) {
	defer func() {
		if err != nil {
			core.Warn("Error obtaining db: %v\n%s", err, string(debug.Stack()))
		}
	}()
	if dsns == nil {
		err = cfgmgr.Load("dsn", &dsns)
		if err != nil {
			core.Warn("No dsn entries found using std values")
			dsns = make(map[string]string)
			dsns["default"] = "host=postgres user=postgres password=postgres dbname=${DBNAME}"
		}

		go func() {
			chcfg := cfgmgr.NotifyChange("dsn")
			for {
				s := <-chcfg
				err = yaml.NewDecoder(strings.NewReader(s)).Decode(&dsns)
				if err != nil {
					core.Err(err)
				}
			}

		}()
	}
	ret, ok := loadDb(n)
	if !ok {
		core.Log("Opening DB: %s", n)
		dsn, ok := dsns[n]
		if !ok {
			dsn, ok = dsns["default"]
			if ok {
				dsn = strings.ReplaceAll(dsn, "${DBNAME}", n)
			} else {
				dsn = "host=postgres user=postgres password=postgres dbname=" + n
			}
		}

		var lerr error

		ret, lerr = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if lerr != nil {
			if strings.Contains(lerr.Error(), "database") && strings.Contains(lerr.Error(), "does not exist") {
				db, err := DBMaster()
				if err != nil {
					return nil, err
				}
				err = db.Exec("create database " + n + ";").Error
				if err != nil {
					return nil, err
				}
			}

			ret, lerr = gorm.Open(postgres.New(postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}), &gorm.Config{})

			if lerr != nil {
				err = lerr
				return
			}
		}

		lerr = ret.Raw("select 1+1").Error
		if lerr != nil {
			err = lerr
			return
		} else {
			err = nil
		}
		core.Debug("DB: New  Connection: %s", n)
		saveDb(n, ret)

	}
	return
}

func DBClose(n string) error {
	db, ok := loadDb(n)
	if ok {
		rdb, err := db.DB()
		if err != nil {
			return err
		}
		err = rdb.Close()
		if err != nil {
			return err
		}
		core.Debug("DB: Closed connection: %s", n)
		delDb(n)
	}
	return nil
}

func DBCloseAll() {
	ks := make([]string, len(dbs))
	for k := range dbs {
		ks = append(ks, k)
	}
	for _, k := range ks {
		DBClose(k)
	}
}

func DSNS() ([]string, error) {
	ks, err := redismgr.Keys("config:dsn:*")
	ret := make([]string, 0)
	if err != nil {
		return nil, err
	}
	for _, k := range ks {
		parts := strings.Split(k, ":")
		ret = append(ret, parts[2])
	}
	return ret, nil
}

func DropRecreate(n string) error {
	db, err := DBMaster()
	if err != nil {
		return err
	}
	err = db.Exec(fmt.Sprintf("drop database %s;", n)).Error
	if err != nil {
		return err
	}
	err = db.Exec(fmt.Sprintf("create database %s;", n)).Error
	if err != nil {
		return err
	}
	return nil
}
