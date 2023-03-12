package srk

import "time"

type Contest struct {
	Title          map[string]string `json:"title"`
	StartAt        time.Time         `json:"startAt"`        // binding:"datetime=2006-01-02T15:04:05Z07:00"`
	EndAt          time.Time         `json:"endAt"`          // binding:"datetime=2006-01-02T15:04:05Z07:00"`
	FrozenDuration time.Duration     `json:"frozenDuration"` // 单位/秒

	Problems []Problem `json:"problems"`
	Members  []Member  `json:"members"`
	Markers  []Marker  `json:"markers"`
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
