package access

import (
	"rankland/database"
	"rankland/model"

	"gorm.io/gorm"
)

func GetRankByID(id int64) (*model.Rank, error) {
	r := &model.Rank{}
	db := database.GetDB().Where("id = ?", id)
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

func GetRankByUniqueKey(uniqueKey string) (*model.Rank, error) {
	r := &model.Rank{}
	db := database.GetDB().Where("unique_key = ?", uniqueKey)
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
	rank := model.Rank{
		UniqueKey: uniqueKey,
		Name:      name,
		Content:   content,
		FileID:    fileID,
	}
	db := database.GetDB()
	if err := db.Create(&rank).Error; err != nil {
		return id, err
	}

	return rank.ID, nil
}

func UpdateRank(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := database.GetDB().Model(&model.Rank{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func SearchRank(query string, pageSize int) (ranks []model.Rank, err error) {
	db := database.GetDB().Where("unique_key like ?", query+"%").Or("name like ?", query+"%")
	if err := db.Limit(pageSize).Find(&ranks).Error; err != nil {
		return ranks, err
	}

	return ranks, nil
}
