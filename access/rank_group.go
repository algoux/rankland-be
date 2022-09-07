package access

import (
	"rankland/database"
	"rankland/model"

	"gorm.io/gorm"
)

func GetRankGroupByID(id int64) (*model.RankGroup, error) {
	rg := &model.RankGroup{}
	db := database.GetDB().Where("id = ?", id)
	if err := db.First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rg, nil
}

func GetRankGroupByName(name string) (*model.RankGroup, error) {
	rg := &model.RankGroup{}
	db := database.GetDB().Where("name = ?", name)
	if err := db.First(&rg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return rg, nil
}

func CreateRankGroup(name, content string) (id int64, err error) {
	rg := &model.RankGroup{
		Name:    name,
		Content: content,
	}
	db := database.GetDB()
	if err := db.Create(rg).Error; err != nil {
		return id, err
	}

	return rg.ID, nil
}

func UpdateRankGroup(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := database.GetDB().Model(&model.RankGroup{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
