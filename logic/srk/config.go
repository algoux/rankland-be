package srk

import (
	"time"
)

type Config struct {
	ID         int64             `json:"id,string"`
	UniqueKey  string            `json:"uniqueKey"`
	Title      map[string]string `json:"title"`
	StartAt    time.Time         `json:"startAt"`        // binding:"datetime=2006-01-02T15:04:05Z07:00"`
	Duration   Duration          `json:"duration"`       // binding:"datetime=2006-01-02T15:04:05Z07:00"`
	Frozen     Duration          `json:"frozenDuration"` // 单位/秒
	UnfrozenAt time.Time         `json:"unfrozenAt,omitempty"`

	Problems     []map[string]any `json:"problems"`
	Members      []map[string]any `json:"members"`
	Markers      []map[string]any `json:"markers,omitempty"`
	Series       []map[string]any `json:"series"`
	Sorter       map[string]any   `json:"sorter"`
	Contributors []string         `json:"contributors"`
	Type         string           `json:"type"`
}

type Problem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Style string `json:"style"`
}

type Member struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Organization string   `json:"organization"`
	TeamMembers  []string `json:"teamMembers"`
	Official     bool     `json:"official"`
	MarkerID     string   `json:"markerID"`
}

type Marker struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Style string `json:"style"`
}
