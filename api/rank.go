package api

import (
	"net/http"
	"ranklist/errcode"
	"ranklist/logic"
	"strconv"

	"github.com/gin-gonic/gin"
)

const OfficialRoot int64 = 1

func GetOfficial(c *gin.Context) {
	id := OfficialRoot
	if strID, ok := c.GetQuery("id"); ok {
		i, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			c.Errors = append(c.Errors, errcode.ParamErr)
			return
		}

		id = i
	}

	data, err := logic.GetRanks(id)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "成功",
		"data":    data,
	})
}

func GetRankNode(c *gin.Context) {
	strID := c.DefaultQuery("id", "1")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	dirs, err := logic.GetChildNodesByID(id)
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

	var uniqueKey string
	if u, ok := any["uniqueKey"]; ok {
		uniqueKey = u.(string)
	} else {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	var parentID int64 = OfficialRoot
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
	var fileID string
	if con, ok := any["fileID"]; ok {
		fileID = con.(string)
	}

	id, err := logic.CreateNode(name, uniqueKey, parentID, int32(typ), fileID)
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
