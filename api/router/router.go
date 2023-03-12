package router

import (
	"net/http"
	"rankland/api"
	"rankland/middleware"

	"github.com/gin-gonic/gin"
)

func Group(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "成功",
			"data":    "hello word",
		})
	})

	r.GET("/statistics", api.GetStatistics)

	rank(r.Group("/rank"))
	contest(r.Group("/contest"))
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

func contest(rg *gin.RouterGroup) {
	rg.GET("/:id", api.GetContest)
	rg.POST("", middleware.WriteHeader(), api.CreateContest)
	rg.PUT("/:id", middleware.WriteHeader(), api.UpdateContest)
	rg.DELETE("/:id", middleware.WriteHeader(), api.DeleteContest)

	rg.POST("/record/:id", middleware.WriteHeader(), api.SetRecord) // 比赛提交记录
	rg.GET("/rank/:id", api.GetRankByContestID)
	rg.GET("/record/:id", api.GetRecordByContestID)
}

func file(rg *gin.RouterGroup) {
	rg.POST("/upload", middleware.WriteHeader(), api.Upload)
	rg.GET("/download", api.Download)
}
