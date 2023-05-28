package api

import (
	"encoding/json"
	"rankland/errcode"
	"rankland/logic"
	"rankland/logic/srk"
	"strconv"

	"github.com/gin-gonic/gin"
)

const adminToken = "sdutpc15_rl_algoux"

func GetRankingConfig(c *gin.Context) {
	key := c.Param("key")
	if id, err := strconv.ParseInt(key, 10, 64); err == nil {
		ct, err := logic.GetRankingConfigByID(id)
		if err != nil {
			c.Errors = append(c.Errors, errcode.ServerErr)
			return
		}

		statusOk(c, ct)
		return
	}

	ct, err := logic.GetRankingByUniqueKey(key)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	statusOk(c, ct)
}

func CreateRankingConfig(c *gin.Context) {
	sc := srk.Config{}
	if err := c.ShouldBindJSON(&sc); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	id, err := logic.CreateRankingConfig(sc)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func UpdateRankingConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	ct := srk.Config{ID: id}
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	err = logic.UpdateRankingConfig(ct)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	statusOk(c, nil)
}

func DeleteContest(c *gin.Context) {
	statusOk(c, nil)
}

func GetRankingByConfigID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	isAdmin := false
	ticket, ok := c.GetQuery("token")
	if ok && ticket == adminToken {
		isAdmin = true
	}
	srkStr, err := logic.GetRankingByConfigID(id, isAdmin)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	srk := make(map[string]interface{})
	err = json.Unmarshal([]byte(srkStr), &srk)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, srk)
}

func GetRecordByConfigID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	isAdmin := false
	ticket, ok := c.GetQuery("token")
	if ok && ticket == adminToken {
		isAdmin = true
	}

	if err := logic.GetRecordByConfigID(id, c.Writer, c.Request, isAdmin); err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
}

func SetRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	records := make([]srk.Record, 0)
	if err := c.ShouldBindJSON(&records); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	err = logic.SetRecord(id, records)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, nil)
}

func ClearRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	err = logic.ClearRecord(id)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, nil)
}
