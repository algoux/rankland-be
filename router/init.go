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

	router.Use(
		middleware.Cors(cors), // 启用跨域拦截
		middleware.Error(),    // 启用 Error 处理
	)

	group(router)

	return router.Run(fmt.Sprintf("%v:%v", host, port))
}

func group(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "成功",
			"data":    "hello word",
		})
	})

	file(r)
	ranklist(r)
}
