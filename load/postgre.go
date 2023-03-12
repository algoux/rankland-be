package load

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitPostgreSQL() error {
	psql := Conf.PostgreSQL
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable TimeZone=%v", psql.Host, psql.Port, psql.Username, psql.Password, psql.DBname, psql.TimeZone)
	postgre, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return err
	}

	if Conf.Application.Env == EnvProd {
		db = postgre
	} else {
		db = postgre.Debug()
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}
