package model

import (
	"rankland/utils"
	"time"

	"gorm.io/gorm"
)

type Contest struct {
	ID       int64    `gorm:"type:bigint;primary_key"`
	Config   string   `gorm:"type:jsonb"`
	Problems []string `gorm:"type:jsonb[]"`
	Members  []string `gorm:"type:jsonb[]"`
	Markers  []string `gorm:"type:jsonb[]"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *Contest) TableName() string {
	return "contest"
}

func (d *Contest) BeforeCreate(tx *gorm.DB) error {
	d.ID = utils.Generator.Generate().Int64()
	return nil
}
