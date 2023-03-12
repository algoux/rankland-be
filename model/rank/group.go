package rank

import (
	"rankland/util"
	"time"

	"gorm.io/gorm"
)

type RankGroup struct {
	ID        int64  `gorm:"type:bigint;primary_key"`
	UniqueKey string `gorm:"type:varchar(200);unique;not null"`
	Name      string `gorm:"type:varchar(200);not null"`
	Content   string `gorm:"type:TEXT"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *RankGroup) TableName() string {
	return "rank_group"
}

func (d *RankGroup) BeforeCreate(tx *gorm.DB) error {
	d.ID = util.Generator.Generate().Int64()
	return nil
}
