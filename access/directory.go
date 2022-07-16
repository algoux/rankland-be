package access

import (
	"ranklist/database"
	"ranklist/model"
)

func GetChildDirsByID(id int64) ([]*model.Directory, error) {
	dirs := []*model.Directory{}
	db := database.GetDB().Where("parent_id = ?", id)
	if err := db.Find(&dirs).Error; err != nil {
		return nil, err
	}

	return dirs, nil
}

func CreateDir(name string, parentID int64, typ int32, content string) (int64, error) {
	dir := &model.Directory{
		Name:     name,
		ParentID: parentID,
		Type:     typ,
		Content:  content,
	}
	if err := database.GetDB().Create(dir).Error; err != nil {
		return 0, err
	}

	return dir.ID, nil
}
