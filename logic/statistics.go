package logic

import "rankland/access"

func GetStatistics() (rankCnt, viewCnt int32, err error) {
	return access.GetRankStatistics()
}
