package access

import (
	"ranklist/database"
	"ranklist/model"
)

func GetAllNodes() (map[int64]model.TreeNode, error) {
	nodes := []model.TreeNode{}
	if err := database.GetDB().Find(&nodes).Error; err != nil {
		return nil, err
	}

	allNode := map[int64]model.TreeNode{}
	for _, node := range nodes {
		allNode[node.ID] = node
	}
	return allNode, nil
}

func GetChildNodesByID(id int64) ([]model.TreeNode, error) {
	nodes := []model.TreeNode{}
	db := database.GetDB().Where("parent_id = ?", id)
	if err := db.Find(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

func CreateNode(name, uniqueKey string, parentID int64, typ int32, content string) (int64, error) {
	node := &model.TreeNode{
		Name:      name,
		UniqueKey: uniqueKey,
		ParentID:  parentID,
		Type:      typ,
		Content:   content,
	}
	if err := database.GetDB().Create(node).Error; err != nil {
		return 0, err
	}

	return node.ID, nil
}
