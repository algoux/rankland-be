package srk

type Record struct {
	ID             int64  `json:"id,string"` // 一个标记提交顺序的唯一 ID
	MemberID       string `json:"memberID"`
	ProblemID      string `json:"problemID"`
	Result         string `json:"result"`
	SubmissionTime int64  `json:"submissionTime"` // 单位：秒
}
