package router

import (
	"ranklist/api"

	"github.com/gin-gonic/gin"
)

func ranklist(r *gin.Engine) {
	rg := r.Group("/rank")
	{
		rg.GET("/dir", api.GetRankNode)
		rg.POST("/dir", api.CreateRankNode)
		rg.GET("/node", api.GetRank)
	}
}
