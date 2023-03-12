package api

import (
	"encoding/json"
	"rankland/errcode"
	"rankland/logic"
	"rankland/logic/srk"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetContest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	contest, err := logic.GetContestByID(id)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	statusOk(c, contest)
}

func CreateContest(c *gin.Context) {
	sc := srk.Contest{}
	if err := c.ShouldBindJSON(&sc); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	id, err := logic.CreateContest(sc)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func UpdateContest(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	ct := srk.Contest{}
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	statusOk(c, nil)
}

func DeleteContest(c *gin.Context) {
	statusOk(c, nil)
}

func GetRankByContestID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	srkStr, err := logic.GetRankByContestID(id)
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

func GetRecordByContestID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	recordStr, err := logic.GetRecordsByContestID(id)
	if err != nil {
		statusOk(c, nil)
		return
	}
	records := []interface{}{}
	err = json.Unmarshal([]byte(recordStr), &records)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, records)
}

func SetRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	records := srk.Records{}
	if err := c.ShouldBindJSON(&records); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	err = logic.SetRecord(id, records.Records)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	statusOk(c, nil)
}
