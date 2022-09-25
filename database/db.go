package database

import (
	"fmt"
	"rankland/model"
	"rankland/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitPostgreSQL() error {
	psql := utils.GetConfig().PostgreSQL
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=%v", psql.Host, psql.Port, psql.Username, psql.Password, psql.DBname, psql.TimeZone)
	postgre, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	db = postgre
	return nil
}

func GetDB() *gorm.DB {
	if utils.GetConfig().Application.Env == utils.EnvDev {
		return db.Debug()
	}
	return db
}

func Migration() error {
	if err := db.AutoMigrate(&model.File{}, &model.Rank{}, &model.RankGroup{}); err != nil {
		return err
	}
	return nil
}
