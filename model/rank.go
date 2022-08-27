package model

import (
	"rankland/utils"
	"time"

	"gorm.io/gorm"
)

type Rank struct {
	ID        int64  `gorm:"primary_key"`
	UniqueKey string `gorm:"type:varchar(200);unique;not null"`
	Name      string `gorm:"type:varchar(200)"`
	Content   string `gorm:"type:MEDIUMTEXT"`
	FileID    int64  `gorm:"column:file_id"`
	ViewCnt   int64  `gorm:"index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *Rank) TableName() string {
	return "rank"
}

func (d *Rank) BeforeCreate(tx *gorm.DB) error {
	d.ID = utils.Generator.Generate().Int64()
	return nil
}
