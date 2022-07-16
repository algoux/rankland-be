package api

import (
	"net/http"
	"ranklist/errcode"
	"ranklist/logic"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRankNode(c *gin.Context) {
	strID := c.DefaultQuery("id", "1")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	dirs, err := logic.GetChildDirsByID(id)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    dirs,
	})
}

func CreateRankNode(c *gin.Context) {
	any := make(logic.Any)
	if err := c.ShouldBindJSON(&any); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	var name string
	if n, ok := any["name"]; ok {
		name = n.(string)
	}
	// 根节点 id = 1
	var parentID int64 = 1
	if pid, ok := any["parentID"]; ok {
		id, err := strconv.ParseInt(pid.(string), 10, 64)
		if err != nil {
			c.Errors = append(c.Errors, errcode.ParamErr)
			return
		}
		parentID = id
	}
	var typ int32
	if t, ok := any["type"]; ok {
		typ = int32(t.(float64))
	}
	var content string
	if con, ok := any["content"]; ok {
		content = con.(string)
	}

	id, err := logic.CreateDir(name, parentID, int32(typ), content)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data": map[string]interface{}{
			"id": strconv.FormatInt(id, 10),
		},
	})
}

// 需要重新规划一下
func GetRank(c *gin.Context) {
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

	_, path, err := logic.GetFileByID(id)
	if err != nil {
		c.Errors = append(c.Errors, errcode.FileReadErr)
		return
	}

	c.File(path)
}
