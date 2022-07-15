package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitSqlite() error {
	s, err := gorm.Open(sqlite.Open("sqlite.db"))
	if err != nil {
		return err
	}

	db = s
	return nil
}

func InitMySQL() error {
	dsn := "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
	m, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return err
	}

	db = m
	return nil
}

func GetDB() *gorm.DB {
	return db
}
