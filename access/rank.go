package access

import (
	"rankland/database"
	"rankland/model"

	"gorm.io/gorm"
)

func GetRankByID(id int64) (rank model.Rank, err error) {
	db := database.GetDB().Where("id = ?", id)
	if err := db.First(&rank).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rank, nil
		}
		return rank, err
	}

	// 查询是浏览量增加 1
	db.Model(&rank).Update("view_cnt", rank.ViewCnt+1)
	return rank, nil
}

func GetRankByUniqueKey(uniqueKey string) (rank model.Rank, err error) {
	db := database.GetDB().Where("unique_key = ?", uniqueKey)
	if err := db.First(&rank).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rank, nil
		}
		return rank, err
	}

	// 查询是浏览量增加 1
	db.Model(&rank).Update("view_cnt", rank.ViewCnt+1)
	return rank, nil
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

func UpdateRank(id int64, uniqueKey, name, content string, fileID int64) error {
	rank := model.Rank{
		ID:        id,
		UniqueKey: uniqueKey,
		Name:      name,
		Content:   content,
		FileID:    fileID,
	}
	db := database.GetDB()
	if err := db.Where("id = ?", id).Updates(&rank).Error; err != nil {
		return err
	}
	return nil
}

func SearchRank(query string) (ranks []model.Rank, err error) {
	db := database.GetDB().Where("unique_key like ?", query+"%").Or("name like ?", query+"%")
	if err := db.Limit(100).Find(&ranks).Error; err != nil {
		return ranks, err
	}

	return ranks, nil
}
