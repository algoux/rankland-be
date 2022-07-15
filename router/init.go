package router

import (
	"fmt"
	"net/http"
	"ranklist/middleware"

	"github.com/gin-gonic/gin"
)

func Init(host, port, cors string) error {
	// 默认开启了 logger 和 recovery
	router := gin.Default()

	// 启用跨域拦截
	router.Use(middleware.Cors(cors))

	group(router)

	return router.Run(fmt.Sprintf("%v:%v", host, port))
}

func group(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
			"data":    "hello word",
		})
	})

	file(r)
}
