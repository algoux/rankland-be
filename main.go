package main

import (
	"fmt"
	"rankland/api/router"
	"rankland/load"
	"rankland/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	load.Init()

	if err := load.InitPostgreSQL(); err != nil {
		logrus.WithError(err).Fatalf("init postgresql failed")
	}
	load.InitRedis()

	app := load.Conf.Application
	// 默认开启了 logger 和 recovery
	r := gin.Default()
	if app.Env == load.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(
		middleware.Cors(app.Cors), // 启用跨域拦截
		middleware.Error(),        // 启用 Error 处理
	)

	router.Group(r)
	r.Run(fmt.Sprintf("%v:%v", app.Host, app.Port))
}
