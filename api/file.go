package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "文件解析错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": map[string]interface{}{
			"filename": file.Filename,
			"fileSize": file.Size,
		},
	})

}

func Download(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
	})
}
