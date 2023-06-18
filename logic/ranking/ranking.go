package ranking

import (
	"encoding/json"
	"time"
)

type Ranking struct {
	ID         int64             `json:"id,string"`
	UniqueKey  string            `json:"uniqueKey"`
	Title      map[string]string `json:"title"`
	StartAt    time.Time         `json:"startAt"`        // binding:"datetime=2006-01-02T15:04:05Z07:00"`
	Duration   time.Duration     `json:"duration"`       // binding:"datetime=2006-01-02T15:04:05Z07:00"`
	Frozen     time.Duration     `json:"frozenDuration"` // 单位/秒
	UnfrozenAt time.Time         `json:"unfrozenAt,omitempty"`

	Problems     []map[string]any `json:"problems"`
	Members      []Member         `json:"members"`
	Markers      []map[string]any `json:"markers,omitempty"`
	Series       []map[string]any `json:"series"`
	Sorter       map[string]any   `json:"sorter"`
	Contributors []string         `json:"contributors"`
	Type         string           `json:"type"`
}

type Problem struct {
	id string

	raw map[string]any
}

func (p Problem) GetID() string {
	return ""
}

type Member struct {
	id string

	raw map[string]any // 原始数据
}

func NewMembers(raws []map[string]any) []Member {
	m := make([]Member, 0, len(raws))
	for _, r := range raws {
		m = append(m, NewMember(r))
	}
	return m
}

func NewMember(raw map[string]any) Member {
	return Member{
		id: raw["id"].(string),

		raw: raw,
	}
}

func (m *Member) ID() string {
	return m.id
}

func (m *Member) ToJson() (string, error) {
	ret, err := json.Marshal(m.raw)
	return string(ret), err
}

type Marker struct {
	raw map[string]any
}

type Serie struct {
	raw map[string]any
}

type Contributors struct {
	raw []string
}
