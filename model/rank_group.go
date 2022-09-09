package model

import (
	"rankland/utils"
	"time"

	"gorm.io/gorm"
)

type RankGroup struct {
	ID        int64  `gorm:"primary_key"`
	UniqueKey string `gorm:"type:varchar(200);unique;not null"`
	Name      string `gorm:"unique;not null"`
	Content   string `gorm:"type:TEXT"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *RankGroup) TableName() string {
	return "ranklist"
}

func (d *RankGroup) BeforeCreate(tx *gorm.DB) error {
	d.ID = utils.Generator.Generate().Int64()
	return nil
}
