package api

import (
	"rankland/errcode"
	"rankland/logic"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func GetRankGroup(c *gin.Context) {
	key := c.Param("key")

	if id, err := strconv.ParseInt(key, 10, 64); err == nil {
		rg, err := logic.GetRankGroupByID(id)
		if err != nil {
			c.Errors = append(c.Errors, errcode.ServerErr)
			return
		}

		if rg != nil {
			statusOk(c, rg)
			return
		}
	}

	rg, err := logic.GetRankGroupByUniqueKey(key)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	if rg == nil {
		c.Errors = append(c.Errors, errcode.NoResultErr)
		return
	}

	statusOk(c, rg)
}

func CreateRankGroup(c *gin.Context) {
	rg := logic.RankGroup{}
	if err := c.ShouldBindJSON(&rg); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	if rg.Name == nil || utf8.RuneCountInString(strings.TrimSpace(*rg.Name)) < 5 || rg.Content == nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	id, err := logic.CreateRankGroup(rg)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func UpdateRankGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	rg := logic.RankGroup{ID: id}
	if err := c.ShouldBindJSON(&rg); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	if rg.Name != nil && utf8.RuneCountInString(strings.TrimSpace(*rg.Name)) < 5 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	rg.ID = id
	if err := logic.UpdateRankGroup(rg); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, nil)
}

func GetRank(c *gin.Context) {
	key := c.Param("key")

	if id, err := strconv.ParseInt(key, 10, 64); err == nil {
		r, err := logic.GetRankByID(id)
		if err != nil {
			c.Errors = append(c.Errors, errcode.ServerErr)
			return
		}

		if r != nil {
			statusOk(c, r)
			return
		}
	}

	r, err := logic.GetRankByUniqueKey(key)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	if r == nil {
		c.Errors = append(c.Errors, errcode.NoResultErr)
		return
	}

	statusOk(c, r)
}

func CreateRank(c *gin.Context) {
	r := logic.Rank{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	if utf8.RuneCountInString(r.UniqueKey) < 5 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	id, err := logic.CreateRank(r)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func UpdateRank(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	r := logic.Rank{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	if (r.Name != nil && utf8.RuneCountInString(strings.TrimSpace(*r.Name)) < 5) || (r.Content == nil && r.FileID == nil) || (r.FileID != nil && *r.FileID < 0) {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	r.ID = id
	if err := logic.UpdateRank(r); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, nil)
}

func SearchRank(c *gin.Context) {
	q, ok := c.GetQuery("query")
	if !ok {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	pageSize := getDefaultQueryInt64(c, "pageSize", 20)

	rs := logic.NewRanks()
	err := rs.Search(q, pageSize)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, rs)
}

func getDefaultQueryInt64(c *gin.Context, key string, defaultVal int) int {
	if v, ok := c.GetQuery(key); ok {
		if val, err := strconv.Atoi(v); err == nil {
			return val
		}
	}

	return defaultVal
}
