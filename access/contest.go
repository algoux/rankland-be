package access

import (
	"rankland/database"
	"rankland/model"

	"gorm.io/gorm"
)

func GetContestByID(id int64) (contest *model.Contest, err error) {
	db := database.GetDB().Where("id = ?", id)
	if err := db.First(&contest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return contest, err
	}

	return contest, nil
}

func CreateContest(contest *model.Contest) (int64, error) {
	if err := database.GetDB().Create(contest).Error; err != nil {
		return 0, err
	}

	return contest.ID, nil
}

func UpdateContest(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := database.GetDB().Model(&model.Contest{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
