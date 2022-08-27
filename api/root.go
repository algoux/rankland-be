package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func success(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code:    200,
		Message: "成功",
		Data:    v,
	})
}