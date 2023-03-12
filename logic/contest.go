package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"rankland/errcode"
	"rankland/load"
	"rankland/logic/srk"
	"rankland/model/contest"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ContestRank   = map[int64]*atomic.Value{}
	ContestRecord = map[int64]*atomic.Value{}
)

func GetContestByID(id int64) (srk.Contest, error) {
	mc, err := contest.GetContestByID(id)
	if err != nil {
		return srk.Contest{}, err
	}

	go func() {
		SetContestRank(id)
	}()
	return transfromSRK(*mc)
}

func CreateContest(ct srk.Contest) (id int64, err error) {
	return contest.Create(srkTransfrom(ct))
}

func transfromSRK(mc contest.Contest) (srk.Contest, error) {
	title := make(map[string]string)
	err := json.Unmarshal([]byte(mc.Title), &title)
	if err != nil {
		return srk.Contest{}, err
	}
	problems := []srk.Problem{}
	err = json.Unmarshal([]byte(mc.Problem), &problems)
	if err != nil {
		return srk.Contest{}, err
	}
	members := []srk.Member{}
	err = json.Unmarshal([]byte(mc.Member), &members)
	if err != nil {
		return srk.Contest{}, err
	}
	markers := []srk.Marker{}
	err = json.Unmarshal([]byte(mc.Marker), &markers)
	if err != nil {
		return srk.Contest{}, err
	}

	return srk.Contest{
		Title:          title,
		StartAt:        mc.StartAt,
		EndAt:          mc.EndAt,
		FrozenDuration: mc.FrozenDuration,
		Problems:       problems,
		Members:        members,
		Markers:        markers,
	}, nil
}

func srkTransfrom(ct srk.Contest) contest.Contest {
	title, _ := json.Marshal(ct.Title)
	problems, _ := json.Marshal(ct.Problems)
	members, _ := json.Marshal(ct.Members)
	markers, _ := json.Marshal(ct.Markers)
	return contest.Contest{
		Title:          string(title),
		StartAt:        ct.StartAt,
		EndAt:          ct.EndAt,
		FrozenDuration: ct.FrozenDuration,
		Problem:        string(problems),
		Member:         string(members),
		Marker:         string(markers),
	}
}

func SetRecord(contestID int64, records []srk.Record) error {
	rMap := make(map[string][]interface{})
	for _, r := range records {
		k := fmt.Sprintf("%v:%v", contestID, r.MemberID)
		v := fmt.Sprintf("%v,%v,%v", r.ProblemID, r.Result, r.SulotionTime.Format(time.RFC3339))
		rMap[k] = append(rMap[k], v)
	}

	ctx := context.Background()
	pipe := load.GetRedis().Pipeline()
	for k, v := range rMap {
		pipe.SAdd(ctx, k, v...)
		pipe.Expire(ctx, k, 30*24*time.Hour)
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}

	go func() {
		SetContestRank(contestID)
		SetContestRecord(contestID, records)
	}()
	return nil
}

func GetRecord(contestID int64, memberIDs []string) (map[string][]srk.Record, error) {
	ctx := context.Background()
	pipe := load.GetRedis().Pipeline()
	for _, id := range memberIDs {
		pipe.SMembers(ctx, fmt.Sprintf("%v:%v", contestID, id))
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	memberRecords := make(map[string][]srk.Record)
	for i, cmd := range cmds {
		rs, err := cmd.(*redis.StringSliceCmd).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		memberID := memberIDs[i]
		records := make([]srk.Record, 0, len(rs))
		for _, r := range rs {
			v := strings.Split(r, ",")
			sulotionTime, err := time.Parse(v[2], time.RFC3339)
			if err != nil {
				return nil, err
			}
			records = append(records, srk.Record{
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

func SetContestRank(id int64) {
	srk, err := GetSRKRank(id)
	if err != nil {
		return
	}
	if v, ok := ContestRank[id]; ok {
		_ = v.Swap(srk)
		return
	}
	v := &atomic.Value{}
	v.Store(srk)
	ContestRank[id] = v
}

func SetContestRecord(id int64, records []srk.Record) {
	rds, _ := json.Marshal(records)
	if v, ok := ContestRecord[id]; !ok {
		_ = v.Swap(string(rds))

	}
	v := &atomic.Value{}
	v.Store(string(rds))
	ContestRecord[id] = v
}

func GetRankByContestID(id int64) (string, error) {
	if cr, ok := ContestRank[id]; !ok || cr.Load().(string) == "" {
		return "", errcode.NoResultErr
	}

	val := ContestRank[id].Load().(string)
	return val, nil
}

func GetRecordsByContestID(id int64) (string, error) {
	if cr, ok := ContestRecord[id]; !ok || cr.Load().(string) == "" {
		return "", errcode.NoResultErr
	}

	val := ContestRecord[id].Load().(string)
	return val, nil
}

func GetSRKRank(contestID int64) (string, error) {
	ct, err := contest.GetContestByID(contestID)
	if err != nil {
		return "", err
	}
	sc, err := transfromSRK(*ct)
	if err != nil {
		return "", err
	}

	memberIDs := make([]string, 0, len(sc.Members))
	for _, m := range sc.Members {
		memberIDs = append(memberIDs, m.ID)
	}
	memberRecords, err := GetRecord(contestID, memberIDs)
	if err != nil {
		return "", err
	}

	srkRank := map[string]interface{}{
		"type":         getType(),
		"version":      getVersion(),
		"sorter":       getSorter(),
		"contributors": getContributors(),
		"series":       getSeries(),
		"contest":      getContest(sc.Title, sc.StartAt, sc.EndAt, sc.FrozenDuration),
		"problems":     getProblems(sc.Problems),
		"rows":         getRows(sc, memberRecords),
		"markers":      getMarkers(sc.Markers),
	}

	v, err := json.Marshal(srkRank)
	return string(v), err
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
				"options": map[string]interface{}{
					"count": map[string][]int64{"value": {10, 20, 30}},
				},
			},
			"segments": []map[string]string{
				{"title": "金奖", "style": "gold"},
				{"title": "银奖", "style": "silver"},
				{"title": "铜奖", "style": "bronze"},
			},
		},
		map[string]interface{}{
			"title": "R#",
			"rule": map[string]interface{}{
				"preset":  "Normal",
				"options": map[string]interface{}{},
			},
		},
		map[string]interface{}{
			"title": "S#",
			"rule": map[string]interface{}{
				"preset": "UniqByUserField",
				"options": map[string]interface{}{
					"field":               "organization",
					"includeOfficialOnly": true,
				},
			},
		},
	}
}

func getContest(title map[string]string, startAt, endAt time.Time, frozen time.Duration) map[string]interface{} {
	return map[string]interface{}{
		"title": map[string]string{
			"zh-CN":    title["zh"],
			"fallback": title["zh"],
			"en-US":    title["en"],
		},
		"startAt":        startAt.Format(time.RFC3339),
		"duration":       []interface{}{endAt.Sub(startAt) / time.Hour, "h"},
		"frozenDuration": []interface{}{frozen, "s"},
	}
}

func getMarkers(markers []srk.Marker) []interface{} {
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

func getProblems(problems []srk.Problem) []interface{} {
	ps := make([]interface{}, 0, len(problems))
	for _, p := range problems {
		ps = append(ps, map[string]interface{}{
			"alias": p.Title,
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
	user     srk.Member
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

func getRows(sc srk.Contest, memberRecords map[string][]srk.Record) []map[string]interface{} {
	rows := make([]row, 0, len(sc.Members))
	for _, member := range sc.Members {
		records, ok := memberRecords[member.ID]
		if !ok {
			rows = append(rows, row{
				allTime:  0,
				value:    0,
				user:     member,
				statuses: make([]map[string]interface{}, len(sc.Problems)),
			})
			continue
		}

		sort.Slice(records, func(i, j int) bool {
			return records[i].SulotionTime.Before(records[j].SulotionTime)
		})
		isSolutions := make(map[string]bool) // 存储题目是否已经被解决
		solutionMap := make(map[string][]solution)
		for _, r := range records {
			if isSolutions[r.ProblemID] {
				continue
			}
			if sc.EndAt.Sub(r.SulotionTime) <= sc.FrozenDuration*time.Second {
				solutionMap[r.ProblemID] = append(solutionMap[r.ProblemID], solution{
					result: "?",
					time:   int64(r.SulotionTime.Sub(sc.StartAt) / time.Second),
				})
				continue
			}

			solutionMap[r.ProblemID] = append(solutionMap[r.ProblemID], solution{
				result: r.Result,
				time:   int64(r.SulotionTime.Sub(sc.StartAt) / time.Second),
			})
			if r.Result == SR_FirstBlood || r.Result == SR_Accepted {
				isSolutions[r.ProblemID] = true
			}
		}

		var allTime, value int64
		stats := make([]map[string]interface{}, len(sc.Problems))
		for i, p := range sc.Problems {
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
				"result":    SR_WrongAnswer,
				"time":      []interface{}{0, "s"},
				"tries":     sLen - 1,
				"solutions": sols,
			}
			if isSolutions[p.ID] {
				stats[i]["result"] = SR_Accepted
				stats[i]["time"] = []interface{}{pTime, "s"}
				allTime += pTime
				value += 1
			}
		}
		rows = append(rows, row{
			allTime:  allTime,
			value:    value,
			user:     member,
			statuses: stats,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].value == rows[j].value {
			return rows[i].allTime > rows[j].allTime
		}
		return rows[i].value < rows[j].value
	})

	rs := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		teamMember := make([]interface{}, 0, len(row.user.TeamMembers))
		for _, t := range row.user.TeamMembers {
			teamMember = append(teamMember, map[string]string{"name": t})
		}
		user := map[string]interface{}{
			"id":           row.user.ID,
			"name":         row.user.Name,
			"organization": row.user.Organization,
			"official":     row.user.Official,
			"teamMembers":  teamMember,
		}
		for _, m := range sc.Markers {
			if row.user.MarkerID == m.ID {
				user["marker"] = m.Label
			}
		}
		rs = append(rs, map[string]interface{}{
			"score":    map[string]interface{}{"value": row.value, "time": []interface{}{row.allTime, "s"}},
			"statuses": row.statuses,
			"user":     user,
		})
	}
	return rs
}
