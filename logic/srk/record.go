package srk

import "time"

type Record struct {
	MemberID     string    `json:"memberID"`
	ProblemID    string    `json:"problemID"`
	Result       string    `json:"result"`
	SulotionTime time.Time `json:"sulotionTime"` // binding:"datetime=2006-01-02T15:04:05Z07:00"`
}

type Records struct {
	Records []Record `json:"records"`
}
