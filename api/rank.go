package api

import (
	"net/http"
	"rankland/errcode"
	"rankland/logic"
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
func GetRankOld(c *gin.Context) {
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

func GetRankGroup(c *gin.Context) {
	key := c.Param("key")
	var err error
	rt := logic.NewRankGroup()

	if rt.ID, err = strconv.ParseInt(key, 10, 64); err == nil {
		if err := rt.GetByID(); err != nil {
			c.Errors = append(c.Errors, errcode.ParamErr)
			return
		}
		success(c, rt)
		return
	}

	rt.Name = key
	if err := rt.GetByName(); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	success(c, rt)
}

func CreateRankGroup(c *gin.Context) {
	rt := logic.NewRankGroup()
	if err := c.ShouldBindJSON(rt); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	if err := rt.Create(); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	success(c, rt)
}

func UpdateRankGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	rt := logic.NewRankGroup()
	if err := c.ShouldBindJSON(rt); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	rt.ID = id
	if err := rt.Update(); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	success(c, rt)
}

func GetRank(c *gin.Context) {
	key := c.Param("key")
	var err error
	r := logic.NewRank()

	if r.ID, err = strconv.ParseInt(key, 10, 64); err == nil {
		if err := r.GetByID(); err != nil {
			c.Errors = append(c.Errors, errcode.ParamErr)
			return
		}
		success(c, r)
		return
	}

	r.UniqueKey = key
	if err := r.GetByUniqueKey(); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	success(c, r)
}

func CreateRank(c *gin.Context) {
	r := logic.NewRank()
	if err := c.ShouldBindJSON(r); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	if err := r.Create(); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	success(c, r)
}

func UpdateRank(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	r := logic.NewRank()
	if err := c.ShouldBindJSON(r); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	r.ID = id
	if err := r.Update(); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	success(c, r)
}

func SearchRank(c *gin.Context) {
	q, ok := c.GetQuery("query")
	if !ok {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	rs := logic.NewRanks()
	err := rs.Search(q)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	success(c, rs)
}
