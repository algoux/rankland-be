package logic

import (
	"rankland/model/rank"
)

func GetStatistics() (rankCnt, viewCnt int32, err error) {
	return rank.GetRankStatistics()
}
