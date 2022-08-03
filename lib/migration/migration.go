package migration

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Mig struct {
	Id   string
	Desc string
	Up   func(db *gorm.DB) error `gorm:"-"`
}

func (m Mig) TableName() string {
	return "_version"
}

func Run(db *gorm.DB, v Mig) error {
	err := db.AutoMigrate(Mig{})
	if err != nil {
		return err
	}
	var existingMig Mig
	err = db.Where("id = ?", v.Id).First(&existingMig).Error
	switch {
	case errors.Is(gorm.ErrRecordNotFound, err):
		logrus.Infof("Migration %s is pending - applying now", v.Id)
		err = v.Up(db)
		if err != nil {
			return err
		}
		err = db.Create(&v).Error
		if err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		logrus.Debugf("Migration %s already applied to db", v.Id)
	}

	return nil
}
