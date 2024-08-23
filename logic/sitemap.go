package logic

import (
	"rankland/model/rank"
)

const (
	SitemapVolCap = 1000
)

func GetRankSitemapVolCount() (volCnt int32, err error) {
	rankCnt, err := rank.GetRankCntStatistics()
	if err != nil {
		return 0, err
	}

	return (rankCnt + SitemapVolCap - 1) / SitemapVolCap, nil
}

func GetRankSitemapVol(volIdx int) (uniqueKeys []string, err error) {
	if volIdx < 1 {
		return []string{}, nil
	}
	offset := (volIdx - 1) * SitemapVolCap
	uniqueKeys, err = rank.GetAllRankUniqueKeys(offset, SitemapVolCap)
	if err != nil {
		return []string{}, err
	}

	return uniqueKeys, nil
}
