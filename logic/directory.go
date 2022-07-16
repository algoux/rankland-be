package logic

import "ranklist/access"

type Any map[string]interface{}

func GetChildDirsByID(id int64) ([]Any, error) {
	dirs, err := access.GetChildDirsByID(id)
	if err != nil {
		return nil, err
	}

	childs := make([]Any, 0, len(dirs))
	for _, dir := range dirs {
		if dir == nil {
			continue
		}

		c := Any{
			"id":       dir.ID,
			"name":     dir.Name,
			"parentID": dir.ParentID,
			"type":     dir.Type,
			"content":  dir.Content,
		}
		childs = append(childs, c)
	}
	return childs, nil
}

func CreateDir(name string, parentID int64, typ int32, content string) (int64, error) {
	return access.CreateDir(name, parentID, typ, content)
}
