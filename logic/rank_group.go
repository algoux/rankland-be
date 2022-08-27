package logic

import (
	"rankland/access"
	"time"
)

type RankGroup struct {
	ID      int64  `json:"id,string"`
	Name    string `json:"name"`
	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewRankGroup() *RankGroup {
	return &RankGroup{}
}

func (rt *RankGroup) GetByID() error {
	RankGroup, err := access.GetRankGroupByID(rt.ID)
	if err != nil {
		return err
	}

	rt.Name = RankGroup.Name
	rt.Content = RankGroup.Content
	return nil
}

func (rt *RankGroup) GetByName() error {
	RankGroup, err := access.GetRankGroupByName(rt.Name)
	if err != nil {
		return err
	}

	rt.ID = RankGroup.ID
	rt.Content = RankGroup.Content
	return nil
}

func (rt *RankGroup) Create() error {
	id, err := access.CreateRankGroup(rt.Name, rt.Content)
	if err != nil {
		return err
	}

	rt.ID = id
	return nil
}

func (rt *RankGroup) Update() error {
	return access.UpdateRankGroup(rt.ID, rt.Name, rt.Content)
}
