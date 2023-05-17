package srk

type Record struct {
	ID        int64  `json:"id"` // 一个标记提交顺序的唯一 ID
	MemberID  string `json:"memberID"`
	ProblemID string `json:"problemID"`
	Result    string `json:"result"`
	Sulotion  int64  `json:"sulotion"` // s
}
