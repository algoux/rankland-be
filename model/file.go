package model

import (
	"ranklist/utils"
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID   int64 `gorm:"primary_key"`
	Name string
	MD5  string
	Path string

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
