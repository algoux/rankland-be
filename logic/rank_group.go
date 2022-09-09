package logic

import (
	"rankland/access"
	"time"
)

type RankGroup struct {
	ID        int64   `json:"id,string"`
	UniqueKey string  `json:"unique_key"`
	Name      *string `json:"name"`
	Content   *string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetRankGroupByID(id int64) (*RankGroup, error) {
	rg, err := access.GetRankGroupByID(id)
	if err != nil {
		return nil, err
	}
	if rg == nil {
		return nil, nil
	}

	return &RankGroup{
		ID:        rg.ID,
		UniqueKey: rg.UniqueKey,
		Name:      &rg.Name,
		Content:   &rg.Content,
		CreatedAt: rg.CreatedAt,
		UpdatedAt: rg.UpdatedAt,
	}, nil
}

func GetRankGroupByUniqueKey(uniqueKey string) (*RankGroup, error) {
	rg, err := access.GetRankGroupByUniqueKey(uniqueKey)
	if err != nil {
		return nil, err
	}
	if rg == nil {
		return nil, nil
	}

	return &RankGroup{
		ID:        rg.ID,
		UniqueKey: rg.UniqueKey,
		Name:      &rg.Name,
		Content:   &rg.Content,
		CreatedAt: rg.CreatedAt,
		UpdatedAt: rg.UpdatedAt,
	}, nil
}

func CreateRankGroup(rg RankGroup) (int64, error) {
	id, err := access.CreateRankGroup(rg.UniqueKey, *rg.Name, *rg.Content)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateRankGroup(rg RankGroup) error {
	updates := make(map[string]interface{})
	if rg.Name != nil {
		updates["name"] = *rg.Name
	}
	if rg.Content != nil {
		updates["content"] = *rg.Content
	}
	return access.UpdateRankGroup(rg.ID, updates)
}
