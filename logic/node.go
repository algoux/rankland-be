package logic

import (
	"encoding/json"
	"ranklist/access"
)

type Any map[string]interface{}

func GetChildNodesByID(id int64) ([]Any, error) {
	dirs, err := access.GetChildNodesByID(id)
	if err != nil {
		return nil, err
	}

	childs := make([]Any, 0, len(dirs))
	for _, dir := range dirs {
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

func CreateNode(name, uniqueKey string, parentID int64, typ int32, fileID string) (int64, error) {
	var val string
	if typ == NodeTypeFile {
		content := map[string]interface{}{"fileID": fileID}
		v, err := json.Marshal(content)
		if err != nil {
			return 0, err
		}
		val = string(v)
	}

	return access.CreateNode(name, uniqueKey, parentID, typ, val)
}
