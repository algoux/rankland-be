package contest

import (
	"rankland/util"
	"time"

	"gorm.io/gorm"
)

type Contest struct {
	ID             int64  `gorm:"type:bigint;primary_key"`
	Title          string `gorm:"type:varchar(500)"`
	StartAt        time.Time
	EndAt          time.Time
	FrozenDuration time.Duration

	Problem string `gorm:"type:text"`
	Member  string `gorm:"type:text"`
	Marker  string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *Contest) TableName() string {
	return "contest"
}

func (d *Contest) BeforeCreate(tx *gorm.DB) error {
	d.ID = util.Generator.Generate().Int64()
	return nil
}
