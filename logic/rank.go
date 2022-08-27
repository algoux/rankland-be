package logic

import (
	"encoding/json"
	"rankland/access"
	"rankland/model"
	"time"
)

const (
	NodeTypeDir  = 1
	NodeTypeFile = 2
)

type RankNode struct {
	ID        int64      `json:"id,string"`
	UniqueKey string     `json:"uniqueKey"`
	Name      string     `json:"name"`
	ParentID  int64      `json:"parentID,string"`
	Type      int32      `json:"type"`
	FileID    string     `json:"fileID,omitempty"`
	Children  []RankNode `json:"children,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetRanks(nodeID int64) (Any, error) {
	allNode, err := access.GetAllNodes()
	if err != nil {
		return nil, err
	}
	if len(allNode) == 0 {
		return nil, nil
	}

	children := map[int64][]int64{}
	for _, node := range allNode {
		if c, ok := children[node.ParentID]; ok {
			children[node.ParentID] = append(c, node.ID)
			continue
		}

		children[node.ParentID] = []int64{node.ID}
	}

	rank, err := getRankNode(nodeID, children, allNode)
	if err != nil {
		return nil, err
	}

	// 如果是官方榜单的根节点，补充部分内容
	if nodeID == 1 {
		rank.ID = nodeID
		rank.Name = "官方榜单"
		rank.Type = NodeTypeDir
		rank.UniqueKey = "official_root"
	}

	return marshal(rank)
}

func getRankNode(id int64, children map[int64][]int64, allNode map[int64]model.TreeNode) (RankNode, error) {
	n := allNode[id]
	node := RankNode{
		ID:        n.ID,
		UniqueKey: n.UniqueKey,
		Name:      n.Name,
		ParentID:  n.ParentID,
		Type:      n.Type,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.Type == NodeTypeFile {
		v := map[string]interface{}{}
		if err := json.Unmarshal([]byte(n.Content), &v); err != nil {
			return RankNode{}, err
		}
		node.FileID = v["fileID"].(string)
		return node, nil
	}

	for _, childID := range children[id] {
		cNode, err := getRankNode(childID, children, allNode)
		if err != nil {
			return RankNode{}, err
		}
		node.Children = append(node.Children, cNode)
	}

	return node, nil
}

func marshal(v interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	val := map[string]interface{}{}
	err = json.Unmarshal(b, &val)
	return val, err
}

type Rank struct {
	ID        int64  `json:"id,string"`
	UniqueKey string `json:"uniqueKey"`
	Name      string `json:"name"`
	Content   string `json:"content,omitempty"`
	FileID    int64  `json:"fileID,string,omitempty"`
	ViewCnt   int64  `json:"viewCnt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewRank() *Rank {
	return &Rank{}
}

func (r *Rank) GetByID() error {
	rank, err := access.GetRankByID(r.ID)
	if err != nil {
		return err
	}

	r.UniqueKey = rank.UniqueKey
	r.Name = rank.Name
	r.Content = rank.Content
	r.FileID = rank.FileID
	r.ViewCnt = rank.ViewCnt
	r.CreatedAt = rank.CreatedAt
	r.UpdatedAt = rank.UpdatedAt
	return nil
}

func (r *Rank) GetByUniqueKey() error {
	rank, err := access.GetRankByUniqueKey(r.UniqueKey)
	if err != nil {
		return err
	}

	r.ID = rank.ID
	r.Name = rank.Name
	r.Content = rank.Content
	r.FileID = rank.FileID
	r.ViewCnt = rank.ViewCnt
	r.CreatedAt = rank.CreatedAt
	r.UpdatedAt = rank.UpdatedAt
	return nil
}

func (r *Rank) Create() error {
	id, err := access.CreateRank(r.UniqueKey, r.Name, r.Content, r.FileID)
	if err != nil {
		return err
	}

	r.ID = id
	return nil
}

func (r *Rank) Update() error {
	return access.UpdateRank(r.ID, r.UniqueKey, r.Name, r.Content, r.FileID)
}

type Ranks struct {
	Ranks []Rank
}

func NewRanks() *Ranks {
	return &Ranks{
		Ranks: []Rank{},
	}
}

func (r *Ranks) Search(query string) error {
	rs, err := access.SearchRank(query)
	if err != nil {
		return err
	}

	for _, rank := range rs {
		r.Ranks = append(r.Ranks, Rank{
			ID:        rank.ID,
			UniqueKey: rank.UniqueKey,
			Name:      rank.Name,
			Content:   rank.Content,
			FileID:    rank.FileID,
			ViewCnt:   rank.ViewCnt,
			CreatedAt: rank.CreatedAt,
			UpdatedAt: rank.UpdatedAt,
		})
	}
	return nil
}
