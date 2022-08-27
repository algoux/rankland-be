package model

import (
	"rankland/utils"
	"time"

	"gorm.io/gorm"
)

// TreeNode 树形层级结构
type TreeNode struct {
	ID        int64 `gorm:"primary_key"`
	UniqueKey string
	Name      string
	ParentID  int64
	Type      int32
	Attribute string
	Content   string
	Extra     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (d *TreeNode) TableName() string {
	return "node_tree"
}

func (d *TreeNode) BeforeCreate(tx *gorm.DB) error {
	d.ID = utils.Generator.Generate().Int64()
	return nil
}
