package migrate

import (
	"ranklist/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Sqlite() {
	db, err := gorm.Open(sqlite.Open("sqlite.db"))
	if err != nil {
		logrus.WithError(err).Infof("init sqlite failed")
		return
	}

	err = db.AutoMigrate(&model.File{}, &model.Directory{})
	if err != nil {
		logrus.WithError(err).Infof("auto migrate failed")
	}
	return
}
