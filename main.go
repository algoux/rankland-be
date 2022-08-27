package main

import (
	"rankland/database"
	"rankland/router"
	"rankland/utils"
)

func main() {
	// database.Sqlite()
	if err := database.InitSqlite(); err != nil {
		panic(err)
	}

	if err := utils.InitGenerator(); err != nil {
		panic(err)
	}

	if err := router.Init("127.0.0.1", "8000", "*"); err != nil {
		panic(err)
	}
}
