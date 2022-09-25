package main

import (
	"rankland/database"
	"rankland/router"
	"rankland/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := utils.InitConfig(); err != nil {
		logrus.WithError(err).Fatalf("init config failed")
	}

	if err := database.InitPostgreSQL(); err != nil {
		logrus.WithError(err).Fatalf("init postgresql failed")
	}

	// DB 数据表迁移
	if utils.GetConfig().Application.Migration {
		if err := database.Migration(); err != nil {
			logrus.WithError(err).Fatalf("migration db table failed")
		}
	}

	if err := router.InitGin(); err != nil {
		logrus.WithError(err).Fatalf("init application failed")
	}
}
