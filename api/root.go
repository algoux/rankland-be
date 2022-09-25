package api

import (
	"net/http"
	"rankland/errcode"
	"rankland/logic"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	rankCnt, viewCnt, err := logic.GetStatistics()
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	statusOk(c, map[string]interface{}{
		"totalSrkCount":  rankCnt,
		"totalViewCount": viewCnt,
	})
}

type Resp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func statusOk(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code:    0,
		Message: "成功",
		Data:    v,
	})
}
