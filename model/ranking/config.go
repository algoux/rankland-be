package ranking

import (
	"rankland/util"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	ID         int64  `gorm:"type:bigint;primary_key"`
	UniqueKey  string `gorm:"type:varchar(200)"`
	Title      string `gorm:"type:varchar(500)"`
	StartAt    time.Time
	EndAt      time.Time
	Frozen     int64
	UnfrozenAt time.Time

	Problem     string `gorm:"type:text"`
	Member      string `gorm:"type:text"`
	Marker      string `gorm:"type:text"`
	Series      string `gorm:"type:text"`
	Sorter      string `gorm:"type:text"`
	Contributor string `gorm:"type:text"`
	Type        string `gorm:"type:varchar(100)"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (c *Config) TableName() string {
	return "ranking_config"
}

func (c *Config) BeforeCreate(tx *gorm.DB) error {
	c.ID = util.Generator.Generate().Int64()
	return nil
}
