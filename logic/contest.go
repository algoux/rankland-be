package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"rankland/access"
	"rankland/database"
	"rankland/interface/contest"
	"rankland/model"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type ContestSeries struct {
	Title    string   `json:"title"`
	Rule     string   `json:"rule"`
	Segments []string `json:"segments"`
}

type ContestConfig struct {
	Title    string          `json:"title"`
	StartAt  time.Time       `json:"startAt"`
	EndAt    time.Time       `json:"endAt"`
	FrozenAt time.Time       `json:"frozenAt"`
	Link     string          `json:"link"`
	Series   []ContestSeries `json:"series"`
	Markers  []string        `json:"markers"`
}

func GetContestByID(id int64) (contest contest.Contest, err error) {
	mc, err := access.GetContestByID(id)
	if err != nil {
		return contest, err
	}

	return transfromContest(*mc)
}

func CreateContest(contest contest.Contest) (id int64, err error) {
	mc := contestTransfrom(contest)
	return access.CreateContest(mc)
}

func transfromContest(mc model.Contest) (contest.Contest, error) {
	config := contest.Config{}
	err := json.Unmarshal([]byte(mc.Config), &config)
	if err != nil {
		return contest.Contest{}, err
	}
	problems := []contest.Problem{}
	err = json.Unmarshal([]byte(mc.Config), &problems)
	if err != nil {
		return contest.Contest{}, err
	}
	members := []contest.Member{}
	err = json.Unmarshal([]byte(mc.Config), &members)
	if err != nil {
		return contest.Contest{}, err
	}
	markers := []contest.Marker{}
	err = json.Unmarshal([]byte(mc.Config), &markers)
	if err != nil {
		return contest.Contest{}, err
	}

	return contest.Contest{
		Config:   config,
		Problems: problems,
		Members:  members,
		Markers:  markers,
	}, nil
}

func contestTransfrom(ct contest.Contest) *model.Contest {
	config, _ := json.Marshal(ct.Config)
	problems := make([]string, 0, len(ct.Problems))
	for _, p := range ct.Problems {
		problem, _ := json.Marshal(p)
		problems = append(problems, string(problem))
	}
	members := make([]string, 0, len(ct.Members))
	for _, m := range ct.Members {
		member, _ := json.Marshal(m)
		members = append(members, string(member))
	}
	markers := make([]string, 0, len(ct.Markers))
	for _, m := range ct.Markers {
		marker, _ := json.Marshal(m)
		markers = append(markers, string(marker))
	}

	return &model.Contest{
		Config:   string(config),
		Problems: problems,
		Members:  members,
		Markers:  markers,
	}
}

var recordChan = make(chan contest.Record, 1000)

func SetRecord(contestID int64, records []contest.Record) error {
	rMap := make(map[string][]interface{})
	for _, r := range records {
		k := fmt.Sprintf("%v:%v", contestID, r.MemberID)
		v := fmt.Sprintf("%v,%v,%v", r.ProblemID, r.Result, r.SulotionTime.Format(time.RFC3339))
		rMap[k] = append(rMap[k], v)
	}

	ctx := context.Background()
	pipe := database.GetRedis().Pipeline()
	for k, v := range rMap {
		pipe.SAdd(ctx, k, v...)
		pipe.Expire(ctx, k, 30*24*time.Hour)
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}
	for _, r := range records {
		recordChan <- r
	}
	return nil
}

func GetRecord(contestID int64, memberIDs []string) (map[string][]contest.Record, error) {
	ctx := context.Background()
	pipe := database.GetRedis().Pipeline()
	for _, id := range memberIDs {
		pipe.SMembers(ctx, fmt.Sprintf("%v:%v", contestID, id))
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	memberRecords := make(map[string][]contest.Record)
	for i, cmd := range cmds {
		rs, err := cmd.(*redis.StringSliceCmd).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		memberID := memberIDs[i]
		records := make([]contest.Record, 0, len(rs))
		for _, r := range rs {
			v := strings.Split(r, ",")
			sulotionTime, err := time.Parse(v[2], time.RFC3339)
			if err != nil {
				return nil, err
			}
			records = append(records, contest.Record{
				MemberID:     memberID,
				ProblemID:    v[0],
				Result:       v[1],
				SulotionTime: sulotionTime,
			})
		}
		memberRecords[memberID] = records
	}
	return memberRecords, nil
}

func GetRankByContestID(contestID int64) (string, error) {
	ct, err := GetContestByID(contestID)
	if err != nil {
		return "", err
	}

}

func getType() string {
	return "general"
}

func getVersion() string {
	return "0.3.0"
}

func getSorter() map[string]interface{} {
	return map[string]interface{}{
		"algorithm": "ICPC",
		"config": map[string]interface{}{
			"penalty": []interface{}{20, "min"},
		},
	}
}

func getContributors() []string {
	return []string{"algoUX (https://algoux.org)"}
}

func getSeries() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"title": "#",
			"rule": map[string]interface{}{
				"preset": "ICPC",
				"option": map[string]interface{}{},
			},
		},
		map[string]interface{}{
			"title": "R#",
			"rule": map[string]interface{}{
				"preset": "Normal",
				"option": map[string]interface{}{},
			},
		},
		map[string]interface{}{
			"title": "S#",
			"rule": map[string]interface{}{
				"preset": "UniqByUserField",
				"option": map[string]interface{}{
					"field":               "organization",
					"includeOfficialOnly": true,
				},
			},
		},
	}
}

func getContest(conf contest.Config) map[string]interface{} {
	return map[string]interface{}{
		"title": map[string]string{
			"zh-CN":    conf.Title["zh"],
			"fallback": conf.Title["zh"],
			"en":       conf.Title["en"],
		},
		"startAt":        conf.StartAt.Format(time.RFC3339),
		"duration":       []interface{}{conf.EndAt.Sub(conf.StartAt) / time.Hour, "h"},
		"frozenDuration": conf.FrozenDuration,
	}
}

func getMarkers(markers []contest.Marker) []interface{} {
	mks := make([]interface{}, 0, len(markers))
	for _, m := range markers {
		mks = append(mks, map[string]interface{}{
			"id":    m.ID,
			"label": m.Label,
			"style": m.Style,
		})
	}
	return mks
}

func getProblems(problems []contest.Problem) []interface{} {
	ps := make([]interface{}, 0, len(problems))
	for _, p := range problems {
		ps = append(ps, map[string]interface{}{
			"alias": p.ID,
			"style": map[string]string{
				"backgroundColor": p.Style,
			},
		})
	}
	return ps
}

type row struct {
	allTime  int64
	value    int64
	user     contest.Member
	statuses []map[string]interface{}
}

type solution struct {
	result string
	time   int64
}

const (
	SR_FirstBlood          = "FB"
	SR_Accepted            = "AC"
	SR_Rejected            = "RJ"
	SR_WrongAnswer         = "WA"
	SR_PresentationError   = "PE"
	SR_TimeLimitExceeded   = "TLE"
	SR_MemoryLimitExceeded = "MLE"
	SR_OutputLimitExceeded = "OLE"
	SR_RuntimeError        = "RTE"
	SR_CompilationError    = "CE"
	SR_UnknownError        = "UKE"
	SR_Frozen              = "?"
)

func getRows(ct contest.Contest, memberRecords map[string][]contest.Record) {
	rows := make([]row, 0, len(ct.Members))
	for _, member := range ct.Members {
		records, ok := memberRecords[member.ID]
		if !ok {
			rows = append(rows, row{
				allTime:  0,
				value:    0,
				user:     member,
				statuses: make([]map[string]interface{}, len(ct.Problems)),
			})
			continue
		}

		sort.Slice(records, func(i, j int) bool {
			return records[i].SulotionTime.Before(records[j].SulotionTime)
		})
		solus := make(map[string]bool)
		solutionMap := make(map[string][]solution)
		for _, r := range records {
			if solus[r.ProblemID] {
				continue
			}
			if ct.Config.EndAt.Sub(r.SulotionTime) <= time.Duration(ct.Config.FrozenDuration)*time.Second {
				solutionMap[r.ProblemID] = append(solutionMap[r.ProblemID], solution{
					result: "?",
					time:   int64(r.SulotionTime.Sub(ct.Config.StartAt) / time.Second),
				})
				continue
			}

			solutionMap[r.ProblemID] = append(solutionMap[r.ProblemID], solution{
				result: r.Result,
				time:   int64(r.SulotionTime.Sub(ct.Config.StartAt) / time.Second),
			})
			if r.Result == SR_FirstBlood || r.Result == SR_Accepted {
				solus[r.ProblemID] = true
			}
		}

		var allTime, value int64
		stats := make([]map[string]interface{}, len(ct.Problems))
		for i, p := range ct.Problems {
			solution, ok := solutionMap[p.ID]
			if !ok {
				stats[i] = map[string]interface{}{
					"result": nil,
					"time":   []interface{}{0, "s"},
					"tries":  0,
				}
				continue
			}

			sLen := len(solution)
			pTime := solution[sLen-1].time + int64(20*60*(sLen-1))

			sols := make([]map[string]interface{}, 0, len(solutionMap[p.ID]))
			for _, s := range solutionMap[p.ID] {
				sols = append(sols, map[string]interface{}{
					"result": s.result,
					"time":   []interface{}{s.time, "s"},
				})
			}
			stats[i] = map[string]interface{}{
				"result":    solutionMap[p.ID][sLen-1].result,
				"time":      []interface{}{},
				"tries":     sLen - 1,
				"solutions": sols,
			}

		}

	}

}
