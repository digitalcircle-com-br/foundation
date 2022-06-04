package dbmgr

import (
	"time"

	"github.com/digitalcircle-com-br/foundation/lib/core"
	"github.com/digitalcircle-com-br/foundation/lib/model"
	"gorm.io/gorm"
)

type DBMigration interface {
	ID() string
	Desc() string
	MatchDB(dsn string) bool
	Run(d *gorm.DB) error
}

type innerMigration struct {
	id     string
	desc   string
	runner func(d *gorm.DB) error
	match  func(d string) bool
}

func (i *innerMigration) ID() string {
	return i.id
}

func (i *innerMigration) Desc() string {
	return i.desc
}

func (i *innerMigration) Run(d *gorm.DB) error {
	return i.runner(d)
}

func (i *innerMigration) MatchDB(dsn string) bool {
	return i.match(dsn)
}

var migrations = make([]DBMigration, 0)

func MigrationAdd(id string, desc string, match func(s string) bool, mig func(d *gorm.DB) error) {
	ret := innerMigration{
		id:     id,
		desc:   desc,
		runner: mig,
		match:  match,
	}
	migrations = append(migrations, &ret)
}

func MigrationRunOnDb(dsn string) error {
	db, err := DBN(dsn)
	if err != nil {
		return err
	}
	for _, mig := range migrations {
		if mig.MatchDB(dsn) {

			err = db.AutoMigrate(&model.DBVersion{})
			if err != nil {
				return err
			}
			ver := &model.DBVersion{}
			err = db.Where("id = ?", mig.ID()).First(ver).Error
			switch {
			case err == gorm.ErrRecordNotFound:
				core.Log("Migration %s not found in db %s - Applying. [%s]", mig.ID(), dsn, mig.Desc())
				err = mig.Run(db)
				if err != nil {
					core.Log("Error applying migration %s - %s", mig.ID(), err.Error())
					return err
				}
				ver.Desc = mig.Desc()
				ver.ID = mig.ID()
				ver.RunAt = time.Now()
				err = db.Create(ver).Error
				if err != nil {
					return err
				}
			case err != nil:
				return err
			default:
				core.Log("Migration already applied at: %s", ver.RunAt.String())
			}
		}
	}

	return nil
}

func MigrationRun() error {
	dsns, err := DSNS()
	if err != nil {
		return err
	}
	for _, dsn := range dsns {
		for _, mig := range migrations {
			if mig.MatchDB(dsn) {

				db, err := DBN(dsn)
				if err != nil {
					return err
				}

				err = db.AutoMigrate(&model.DBVersion{})
				if err != nil {
					return err
				}
				ver := &model.DBVersion{}
				err = db.Where("id = ?", mig.ID()).First(ver).Error
				switch {
				case err == gorm.ErrRecordNotFound:
					core.Log("Migration %s not found in db %s - Applying. [%s]", mig.ID(), dsn, mig.Desc())
					err = mig.Run(db)
					if err != nil {
						core.Log("Error applying migration %s - %s", mig.ID(), err.Error())
						return err
					}
					ver.Desc = mig.Desc()
					ver.ID = mig.ID()
					ver.RunAt = time.Now()
					err = db.Create(ver).Error
					if err != nil {
						return err
					}
				case err != nil:
					return err
				default:
					core.Log("Migration already applied at: %s", ver.RunAt.String())
				}
			}
		}

	}
	return nil
}
