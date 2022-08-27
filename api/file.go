package api

import (
	"bytes"
	"net/http"
	"rankland/errcode"
	"rankland/logic"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Errors = append(c.Errors, errcode.FileAnalysisErr)
		return
	}

	fbuf := bytes.Buffer{}
	f, err := file.Open()
	if err != nil {
		c.Errors = append(c.Errors, errcode.FileAnalysisErr)
		return
	}
	if _, err := fbuf.ReadFrom(f); err != nil {
		c.Errors = append(c.Errors, errcode.FileAnalysisErr)
		return
	}

	id, err := logic.CreateFile(file.Filename, fbuf.Bytes())
	if err != nil {
		c.Errors = append(c.Errors, errcode.FileWriteErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data": map[string]interface{}{
			"id": strconv.FormatInt(id, 10),
		},
	})
}

func Download(c *gin.Context) {
	strID, ok := c.GetQuery("id")
	if !ok {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	name, path, err := logic.GetFileByID(id)
	if err != nil {
		c.Errors = append(c.Errors, errcode.FileReadErr)
		return
	}

	// c.File(name)
	c.FileAttachment(path, name)
}
