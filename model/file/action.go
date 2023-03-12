package file

import (
	"rankland/load"

	"gorm.io/gorm"
)

func GetFileByID(id int64) (file *File, err error) {
	db := load.GetDB().Where("id = ?", id)
	if err := db.First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return file, err
	}

	return file, nil
}

func GetFileID(md5 string) (int64, error) {
	file := &File{}
	db := load.GetDB().Select("id").Where("md5 = ?", md5)
	if err := db.First(file).Error; err != nil {
		return 0, err
	}

	return file.ID, nil
}

func CreateFile(name, md5, path string) (int64, error) {
	file := &File{
		Name:   name,
		Secret: md5,
		Path:   path,
	}
	if err := load.GetDB().Create(file).Error; err != nil {
		return 0, err
	}

	return file.ID, nil
}
