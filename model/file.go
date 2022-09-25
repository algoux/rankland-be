package model

import (
	"rankland/utils"
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID     int64  `gorm:"type:bigint;primary_key"`
	Name   string `gorm:"type:varchar(200);not null"`
	Secret string `gorm:"type:varchar(200);not null"`
	Path   string `gorm:"type:varchar(500);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (f *File) TableName() string {
	return "file"
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	f.ID = utils.Generator.Generate().Int64()
	return nil
}
