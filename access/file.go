package access

import (
	"ranklist/database"
	"ranklist/model"
)

func GetFileByID(id int64) (*model.File, error) {
	file := &model.File{}
	db := database.GetDB().Where("id = ?", id)
	if err := db.First(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func GetFileID(md5 string) (int64, error) {
	file := &model.File{}
	db := database.GetDB().Select("id").Where("md5 = ?", md5)
	if err := db.First(file).Error; err != nil {
		return 0, err
	}

	return file.ID, nil
}

func CreateFile(name, md5, path string) (int64, error) {
	file := &model.File{
		Name: name,
		MD5:  md5,
		Path: path,
	}
	if err := database.GetDB().Create(file).Error; err != nil {
		return 0, err
	}

	return file.ID, nil
}
