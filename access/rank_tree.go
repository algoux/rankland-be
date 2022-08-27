package access

import (
	"rankland/database"
	"rankland/model"

	"gorm.io/gorm"
)

func GetRankGroupByID(id int64) (rt model.RankGroup, err error) {
	db := database.GetDB().Where("id = ?", id)
	if err := db.First(&rt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rt, nil
		}
		return rt, err
	}

	return rt, nil
}

func GetRankGroupByName(name string) (rt model.RankGroup, err error) {
	db := database.GetDB().Where("name = ?", name)
	if err := db.First(&rt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return rt, nil
		}
		return rt, err
	}

	return rt, nil
}

func CreateRankGroup(name, content string) (id int64, err error) {
	rt := model.RankGroup{
		Name:    name,
		Content: content,
	}
	db := database.GetDB()
	if err := db.Create(&rt).Error; err != nil {
		return id, err
	}

	return rt.ID, nil
}

func UpdateRankGroup(id int64, name, content string) error {
	rank := model.Rank{
		ID:      id,
		Name:    name,
		Content: content,
	}
	db := database.GetDB()
	if err := db.Where("id = ?", id).Updates(&rank).Error; err != nil {
		return err
	}
	return nil
}
