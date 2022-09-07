package logic

import (
	"rankland/access"
	"time"
)

type Rank struct {
	ID        int64   `json:"id,string"`
	UniqueKey string  `json:"uniqueKey"`
	Name      *string `json:"name"`
	Content   *string `json:"content,omitempty"`
	FileID    *int64  `json:"fileID,string,omitempty"`
	ViewCnt   int64   `json:"viewCnt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetRankByID(id int64) (*Rank, error) {
	r, err := access.GetRankByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}

	return &Rank{
		ID:        r.ID,
		UniqueKey: r.UniqueKey,
		Name:      &r.Name,
		Content:   &r.Content,
		FileID:    &r.FileID,
		ViewCnt:   r.ViewCnt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func GetRankByUniqueKey(uniqueKey string) (*Rank, error) {
	r, err := access.GetRankByUniqueKey(uniqueKey)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}

	return &Rank{
		ID:        r.ID,
		UniqueKey: r.UniqueKey,
		Name:      &r.Name,
		Content:   &r.Content,
		FileID:    &r.FileID,
		ViewCnt:   r.ViewCnt,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func CreateRank(r Rank) (int64, error) {
	var name, content string
	if r.Name != nil {
		name = *r.Name
	}
	if r.Content != nil {
		content = *r.Content
	}
	var fileID int64
	if r.FileID != nil {
		fileID = *r.FileID
	}

	id, err := access.CreateRank(r.UniqueKey, name, content, fileID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateRank(r Rank) error {
	updates := make(map[string]interface{})
	if r.Name != nil {
		updates["name"] = *r.Name
	}
	if r.Content != nil {
		updates["content"] = *r.Content
	}
	if r.FileID != nil {
		updates["file_id"] = *r.FileID
	}
	return access.UpdateRank(r.ID, updates)
}

type Ranks struct {
	Ranks     []Rank `json:"ranks"`
	PageToken int64  `json:"pageToken,string,omitempty"`
}

func NewRanks() *Ranks {
	return &Ranks{
		Ranks: []Rank{},
	}
}

func (r *Ranks) Search(query string, pageSize int) error {
	rs, err := access.SearchRank(query, pageSize)
	if err != nil {
		return err
	}

	for _, rank := range rs {
		r.Ranks = append(r.Ranks, Rank{
			ID:        rank.ID,
			UniqueKey: rank.UniqueKey,
			Name:      &rank.Name,
			Content:   &rank.Content,
			FileID:    &rank.FileID,
			ViewCnt:   rank.ViewCnt,
			CreatedAt: rank.CreatedAt,
			UpdatedAt: rank.UpdatedAt,
		})
	}
	if len(rs) == pageSize {
		r.PageToken = rs[pageSize-1].ID
	}
	return nil
}
