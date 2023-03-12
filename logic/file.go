package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"rankland/model/file"
	"time"
)

const fileDir = "file"

func GetFileByID(id int64) (string, string, error) {
	f, err := file.GetFileByID(id)
	if err != nil {
		return "", "", err
	}
	if f == nil {
		return "", "", nil
	}

	return f.Name, fmt.Sprintf("%v/%v/%v", fileDir, f.Path, f.Name), nil
}

func GetFileID(md5 string) (int64, error) {
	return file.GetFileID(md5)
}

func CreateFile(name string, f []byte) (int64, error) {
	md5 := getFileMD5(f)
	path := getFilePath()

	// 此处先进行保存文件，再将文件信息写入 db，是为了防止数据写入 db 后，文件保存出错产生 db 脏数据
	if err := writeFile(name, path, f); err != nil {
		return 0, err
	}

	id, err := file.CreateFile(name, md5, path)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func getFilePath() string {
	t := time.Now()

	return fmt.Sprintf("%v%02d", t.Year(), t.Month())
}

func getFileMD5(file []byte) string {
	hash := md5.Sum(file)
	return hex.EncodeToString(hash[:])
}

func readFile(name, path string) ([]byte, error) {
	dn := fmt.Sprintf("%v/%v/%v", fileDir, path, name)
	return os.ReadFile(dn)
}

func writeFile(name, path string, file []byte) error {
	err := os.MkdirAll(fmt.Sprintf("%v/%v", fileDir, path), 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}
	dn := fmt.Sprintf("%v/%v/%v", fileDir, path, name)
	return os.WriteFile(dn, file, 0666)
}

type File struct {
	ID     int64  `gorm:"primary_key"`
	Name   string `gorm:"type:varchar(200)"`
	Secret string `gorm:"type:varchar(200)"`
	Path   string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewFile() *File {
	return &File{}
}

func (f *File) GetByID() error {
	file, err := file.GetFileByID(f.ID)
	if err != nil {
		return err
	}

	f.Name = file.Name
	return nil
}
