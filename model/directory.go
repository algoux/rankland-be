package model

import (
	"ranklist/utils"
	"time"

	"gorm.io/gorm"
)

type Directory struct {
	ID       int64 `gorm:"primary_key"`
	Name     string
	ParentID string
	Type     string
	Content  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *Directory) TableName() string {
	return "directory"
}

func (d *Directory) BeforeCreate(tx *gorm.DB) error {
	d.ID = utils.Generator.Generate().Int64()
	return nil
}
