package rank

import (
	"rankland/load"

	"gorm.io/gorm"
)

func GetRankByID(id int64) (*Rank, error) {
	r := &Rank{}
	db := load.GetDB().Where("id = ?", id)
	if err := db.First(r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 查询是浏览量增加 1
	r.ViewCnt++
	db.Model(r).Update("view_cnt", r.ViewCnt)
	return r, nil
}

func GetRankByUniqueKey(uniqueKey string) (*Rank, error) {
	r := &Rank{}
	db := load.GetDB().Where("unique_key = ?", uniqueKey)
	if err := db.First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 查询是浏览量增加 1
	r.ViewCnt++
	db.Model(&r).Update("view_cnt", r.ViewCnt)
	return r, nil
}

func CreateRank(uniqueKey, name, content string, fileID int64) (id int64, err error) {
	rank := Rank{
		UniqueKey: uniqueKey,
		Name:      name,
		Content:   content,
		FileID:    fileID,
	}
	db := load.GetDB()
	if err := db.Create(&rank).Error; err != nil {
		return id, err
	}

	return rank.ID, nil
}

func UpdateRank(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := load.GetDB().Model(&Rank{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func SearchRank(query string, pageSize int) (ranks []Rank, err error) {
	db := load.GetDB().Where("to_tsvector('zh', unique_key || ' ' || name) @@ to_tsquery('zh', ?) AND deleted_at IS NULL", query)
	if err := db.Limit(pageSize).Find(&ranks).Error; err != nil {
		return ranks, err
	}

	return ranks, nil
}

type RankStatistics struct {
	RankCnt int32
	ViewCnt int32
}

func GetRankStatistics() (rankCnt, ViewCnt int32, err error) {
	rs := RankStatistics{}
	db := load.GetDB().Model(&Rank{})
	sql := db.Select("count(*) as rank_cnt", "sum(view_cnt) as view_cnt").Find(&rs)
	if sql.Error != nil {
		return 0, 0, sql.Error
	}

	return rs.RankCnt, rs.ViewCnt, nil
}

func GetRankGroupByID(id int64) (*RankGroup, error) {
	rg := &RankGroup{}
	db := load.GetDB().Where("id = ?", id)
	if err := db.First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rg, nil
}

func GetRankGroupByUniqueKey(uniqueKey string) (*RankGroup, error) {
	rg := &RankGroup{}
	db := load.GetDB().Where("unique_key = ?", uniqueKey)
	if err := db.First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rg, nil
}

func CreateRankGroup(uniqueKey, name, content string) (id int64, err error) {
	rg := &RankGroup{
		UniqueKey: uniqueKey,
		Name:      name,
		Content:   content,
	}
	db := load.GetDB()
	if err := db.Create(rg).Error; err != nil {
		return id, err
	}

	return rg.ID, nil
}

func UpdateRankGroup(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := load.GetDB().Model(&RankGroup{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func ListAllRank() (ranks []Rank, err error) {
	db := load.GetDB().Where("deleted_at IS NULL").Order("created_at DESC")
	if err := db.Find(&ranks).Error; err != nil {
		return ranks, err
	}

	return ranks, nil
}
