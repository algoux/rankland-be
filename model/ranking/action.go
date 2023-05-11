package ranking

import (
	"rankland/load"

	"gorm.io/gorm"
)

func GetConfigByID(id int64) (cfg *Config, err error) {
	db := load.GetDB().Where("id = ?", id)
	if err := db.First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return cfg, err
	}

	return cfg, nil
}

func GetConfigByUniqueKey(key string) (cfg *Config, err error) {
	db := load.GetDB().Where("unique_key = ?", key)
	if err := db.First(&cfg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return cfg, err
	}

	return cfg, nil
}

func Create(cfg Config) (int64, error) {
	if err := load.GetDB().Create(&cfg).Error; err != nil {
		return 0, err
	}

	return cfg.ID, nil
}

func Update(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	db := load.GetDB().Model(&Config{})
	if err := db.Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}
