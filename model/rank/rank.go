package rank

import (
	"rankland/util"
	"time"

	"gorm.io/gorm"
)

type Rank struct {
	ID        int64  `gorm:"type:bigint;primary_key"`
	UniqueKey string `gorm:"type:varchar(200);unique;not null"`
	Name      string `gorm:"type:varchar(200);not null"`
	Content   string `gorm:"type:TEXT"`
	FileID    int64  `gorm:"type:bigint;column:file_id"`
	ViewCnt   int64  `gorm:"type:int;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *Rank) TableName() string {
	return "rank"
}

func (d *Rank) BeforeCreate(tx *gorm.DB) error {
	d.ID = util.Generator.Generate().Int64()
	return nil
}
