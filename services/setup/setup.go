package setup

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/digitalcircle-com-br/foundation/lib/redismgr"
	"github.com/digitalcircle-com-br/foundation/services/auth/hash"
	"github.com/sirupsen/logrus"
)

func loadRedis() error {
	return filepath.WalkDir("keys", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			logrus.Infof("%s - %s", path, d.Name())
			kname := strings.Replace(path, "keys/", "", 1)
			kname = strings.Replace(kname, "/", ":", -1)
			ext := filepath.Ext(kname)
			if ext != "" {
				kname = strings.Replace(kname, ext, "", 1)
			}
			kval, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			redismgr.Set(kname, string(kval))
		}

		return nil
	})
}

func cleanRedis() error {
	return filepath.WalkDir("keys", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			logrus.Infof("%s - %s", path, d.Name())
			kname := strings.Replace(path, "keys/", "", 1)
			kname = strings.Replace(kname, "/", ":", -1)
			redismgr.Del(kname)
		}

		return nil
	})
}

func Clean() error {
	return cleanRedis()
}

func Load() error {
	return loadRedis()
}

// func prepareMigrations() {
// 	dbmgr.MigrationAdd("auth-00001", "Creates Authentication DB",
// 		func(s string) bool {
// 			return s == "auth"
// 		},
// 		func(adb *gorm.DB) error {
// 			for _, mod := range []interface{}{
// 				&model.SecUser{},
// 				&model.SecGroup{},
// 				&model.SecPerm{},
// 			} {
// 				err := adb.AutoMigrate(mod)
// 				if err != nil {
// 					return err
// 				}
// 			}

// 			perm := &model.SecPerm{Name: "*", Val: "*"}
// 			err := adb.Create(perm).Error
// 			if err != nil {
// 				return err
// 			}
// 			group := &model.SecGroup{Name: "root", Perms: []*model.SecPerm{perm}}
// 			err = adb.Create(group).Error
// 			if err != nil {
// 				return err
// 			}
// 			enabled := true

// 			user := &model.SecUser{
// 				Username: "root",
// 				Hash:     "$argon2id$v=19$m=65536,t=3,p=2$nTPFgXmlMFphn506a/VQ2Q$0Y/KXMMxDb28CzuqGZdShAnNuNs3l3vInJRh3xd5uq4",
// 				Email:    "root@root.com",
// 				Enabled:  &enabled, Groups: []*model.SecGroup{group},
// 			}

// 			err = adb.Create(user).Error
// 			if err != nil {
// 				return err
// 			}
// 			return nil

// 		})

// }

// func Run() error {
// 	core.Init("setup")
// 	dsns, err := dbmgr.DSNS()
// 	if err != nil {
// 		return err
// 	}
// 	mdb, err := dbmgr.DBN("postgres")
// 	if err != nil {
// 		return err
// 	}
// 	for _, dsn := range dsns {
// 		if dsn != "postgres" {
// 			err = mdb.Exec(fmt.Sprintf("CREATE DATABASE %s;", dsn)).Error
// 			if err != nil {
// 				log.Printf("DB Creation: %s", err.Error())
// 			}
// 		}
// 	}

// 	prepareMigrations()
// 	dbmgr.MigrationRun()

// 	return nil
// }

// func Drop() error {
// 	core.Init("setup")
// 	dsns, err := dbmgr.DSNS()
// 	if err != nil {
// 		return err
// 	}
// 	mdb, err := dbmgr.DBN("postgres")
// 	if err != nil {
// 		log.Printf("could not connect to master dbmgr.")
// 		return nil
// 	}
// 	for _, dsn := range dsns {
// 		if dsn != "postgres" {

// 			err = mdb.Exec(fmt.Sprintf("DROP DATABASE %s;", dsn)).Error
// 			if err != nil {
// 				log.Printf("DB Creation: %s", err.Error())
// 			}
// 		}
// 	}

// 	return nil
// }

func CreatePassHash(in string) (string, error) {
	return hash.Hash(in)
}
