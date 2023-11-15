package api

import (
	"fmt"
	"net/http"
	"rankland/errcode"
	"rankland/logic"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	rankCnt, viewCnt, err := logic.GetStatistics()
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	statusOk(c, map[string]interface{}{
		"totalSrkCount":  rankCnt,
		"totalViewCount": viewCnt,
	})
}

func SitemapRanklistIndex(c *gin.Context) {
	volCnt, err := logic.GetRankSitemapVolCount()
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	vols := make([]string, 0, volCnt)
	for i := int32(1); i <= volCnt; i++ {
		vols = append(vols, fmt.Sprintf("https://rl.algoux.org/sitemap/ranklist_vol_%d.txt", i))
	}

	c.XML(http.StatusOK, sitemapindex{
		Xmlns:  "http://www.sitemaps.org/schemas/sitemap/0.9",
		Values: vols,
	})
}

func SitemapRanklistVol(c *gin.Context) {
	volName := c.Param("volName")
	idx, err := strconv.ParseInt(strings.TrimSuffix(volName, ".txt"), 10, 32)
	if err != nil || idx <= 0 {
		c.Errors = append(c.Errors, errcode.ParamErr)
		return
	}

	uniqueKeys, err := logic.GetRankSitemapVol(int(idx))
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}

	text := ""
	for _, key := range uniqueKeys {
		text += fmt.Sprintf("https://rl.algoux.org/ranklist/%s\n", key)
	}

	c.String(http.StatusOK, text)
}

type Resp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func statusOk(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, Resp{
		Code:    0,
		Message: "成功",
		Data:    v,
	})
}

type sitemapindex struct {
	Xmlns  string   `xml:"xmlns,attr"`
	Values []string `xml:"sitemap>loc"`
}
