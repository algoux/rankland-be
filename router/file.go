package router

import (
	"ranklist/api"

	"github.com/gin-gonic/gin"
)

func file(r *gin.Engine) {
	rg := r.Group("/file")
	{
		rg.POST("/upload", api.Upload)
		rg.GET("/download", api.Download)
	}
}
