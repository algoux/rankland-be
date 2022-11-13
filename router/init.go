package router

import (
	"fmt"
	"net/http"
	"rankland/api"
	"rankland/middleware"
	"rankland/utils"

	"github.com/gin-gonic/gin"
)

func InitGin() error {
	app := utils.GetConfig().Application
	// 默认开启了 logger 和 recovery
	router := gin.Default()
	if app.Env == utils.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(
		middleware.Cors(app.Cors), // 启用跨域拦截
		middleware.Error(),        // 启用 Error 处理
	)

	group(router)
	return router.Run(fmt.Sprintf("%v:%v", app.Host, app.Port))
}

func group(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "成功",
			"data":    "hello word",
		})
	})

	r.GET("/statistics", api.GetStatistics)

	rank(r.Group("/rank"))
	file(r.Group("/file"))
}

func rank(rg *gin.RouterGroup) {
	rg.POST("", middleware.WriteHeader(), api.CreateRank)
	rg.GET("/:key", api.GetRank)
	rg.PUT("/:id", middleware.WriteHeader(), api.UpdateRank)
	rg.GET("/search", api.SearchRank)

	g := rg.Group("/group")
	g.POST("", middleware.WriteHeader(), api.CreateRankGroup)
	g.GET("/:key", api.GetRankGroup)
	g.PUT("/:id", middleware.WriteHeader(), api.UpdateRankGroup)

}

func file(rg *gin.RouterGroup) {
	rg.POST("/upload", middleware.WriteHeader(), api.Upload)
	rg.GET("/download", api.Download)
}
