package router

import (
	"fmt"
	"net/http"
	"rankland/api"
	"rankland/middleware"

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

	rank(r.Group("/rank"))
	file(r.Group("/file"))
	ranklist(r)
}

func rank(rg *gin.RouterGroup) {
	rg.GET("/group/:key", api.GetRankGroup)
	rg.POST("/group", api.CreateRankGroup)
	rg.PUT("/group/:id", api.UpdateRankGroup)

	rg.GET("/:key", api.GetRank)
	rg.POST("/", api.CreateRank)
	rg.PUT("/:id", api.UpdateRank)
	rg.GET("/search", api.SearchRank)
}

func file(rg *gin.RouterGroup) {
	rg.POST("/upload", api.Upload)
	rg.GET("/download", api.Download)
}

func ranklist(r *gin.Engine) {
	rg := r.Group("/rank")
	{
		// rg.GET("/dir", api.GetRankNode)
		// rg.GET("/node", api.GetRank)
		rg.POST("/official", api.CreateRankNode)
		rg.GET("/official", api.GetOfficial)
	}
}
