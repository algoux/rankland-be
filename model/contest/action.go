package contest

import (
	"rankland/load"

	"gorm.io/gorm"
)

func GetContestByID(id int64) (contest *Contest, err error) {
	db := load.GetDB().Where("id = ?", id)
	if err := db.First(&contest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return contest, err
	}

	return contest, nil
}

func Create(contest Contest) (int64, error) {
	if err := load.GetDB().Create(&contest).Error; err != nil {
		return 0, err
	}

	return contest.ID, nil
}

func Update(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := load.GetDB().Model(&Contest{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
