package api

import (
	"rankland/errcode"
	"rankland/interface/contest"
	"rankland/logic"
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
	ct := contest.Contest{}
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	id, err := logic.CreateContest(ct)
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

	ct := contest.Contest{}
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	statusOk(c, nil)
}

func DeleteContest(c *gin.Context) {

}

func getRankByContestID(c *gin.Context) {

}

func SetRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}
	records := contest.Records{}
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
