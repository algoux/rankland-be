package logic

import (
	"rankland/model/rank"
	"sort"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

func TestRanks_Search(t *testing.T) {
	defer gomonkey.ApplyFunc(rank.SearchRank, func(query string, pageSize int) (ranks []rank.Rank, err error) {
		return []rank.Rank{
			{ID: 1, UniqueKey: "1", Name: "1", Content: "1", ViewCnt: 1},
			{ID: 2, UniqueKey: "2", Name: "2", Content: "2", ViewCnt: 2},
			{ID: 3, UniqueKey: "3", Name: "3", Content: "3", ViewCnt: 3},
		}, nil
	}).Reset()

	ranks := NewRanks()
	ranks.Search("1", 20)
	for i, rank := range ranks.Ranks {
		t.Logf("%v, %+v", i, rank)
	}
}

func TestXxx(t *testing.T) {
	sss := []int64{3, 2, 7, 9}

	sort.Slice(sss, func(i, j int) bool {
		return sss[i] < sss[j]
	})
	t.Errorf("%+v", sss)
}
