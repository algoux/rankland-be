package database

import (
	"rankland/model"

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

	err = db.AutoMigrate(&model.File{}, &model.Rank{}, &model.RankGroup{}, &model.TreeNode{})
	if err != nil {
		logrus.WithError(err).Infof("auto migrate failed")
	}
}
