package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rankland/errcode"
	"rankland/load"
	"rankland/logic/pubsub"
	"rankland/logic/srk"
	"rankland/logic/ws"
	"rankland/model/ranking"
	"rankland/util"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

const DateTime = 1672502400 // 2023-01-01 00:00:00

var (
	Ranking = map[int64]*atomic.Value{}
)

func GetRankingConfigByID(id int64) (srk.Config, error) {
	mc, err := ranking.GetConfigByID(id)
	if err != nil {
		return srk.Config{}, err
	}
	if mc == nil {
		return srk.Config{}, nil
	}

	return transfromSRK(*mc)
}

func GetRankingByUniqueKey(uk string) (srk.Config, error) {
	mc, err := ranking.GetConfigByUniqueKey(uk)
	if err != nil {
		return srk.Config{}, err
	}
	if mc == nil {
		return srk.Config{}, nil
	}

	return transfromSRK(*mc)
}

func CreateRankingConfig(ct srk.Config) (id int64, err error) {
	c, err := srkTransfrom(ct)
	if err != nil {
		return 0, err
	}
	return ranking.Create(c)
}

func UpdateRankingConfig(ct srk.Config) error {
	c, err := srkTransfrom(ct)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if strings.Trim(c.UniqueKey, " ") != "" {
		updates["unique_key"] = strings.Trim(c.UniqueKey, " ")
	}
	if c.Title != "null" && strings.Trim(c.Title, " ") != "" {
		updates["title"] = strings.Trim(c.Title, " ")
	}
	if !c.StartAt.IsZero() && c.StartAt.After(time.Unix(DateTime, 0)) {
		updates["start_at"] = c.StartAt
	}
	if !c.EndAt.IsZero() && c.EndAt.After(time.Unix(DateTime, 0)) && c.EndAt.After(c.StartAt) {
		updates["end_at"] = c.EndAt
	}
	if c.Frozen > 0 {
		updates["frozen"] = c.Frozen
	}
	if !c.UnfrozenAt.IsZero() && c.UnfrozenAt.After(time.Unix(DateTime, 0)) && c.UnfrozenAt.After(c.StartAt) {
		updates["unfrozen_at"] = c.UnfrozenAt
	}

	if c.Problem != "null" && strings.Trim(c.Problem, " ") != "" {
		updates["problem"] = strings.Trim(c.Problem, " ")
	}
	if c.Member != "null" && strings.Trim(c.Member, " ") != "" {
		updates["member"] = strings.Trim(c.Member, " ")
	}
	if c.Marker != "null" && strings.Trim(c.Marker, " ") != "" {
		updates["marker"] = strings.Trim(c.Marker, " ")
	}
	if c.Series != "null" && strings.Trim(c.Series, " ") != "" {
		updates["series"] = strings.Trim(c.Series, " ")
	}
	if c.Sorter != "null" && strings.Trim(c.Sorter, " ") != "" {
		updates["sorter"] = strings.Trim(c.Sorter, " ")
	}
	if c.Contributor != "null" && strings.Trim(c.Contributor, " ") != "" {
		updates["contributor"] = strings.Trim(c.Contributor, " ")
	}
	if c.Type != "null" && strings.Trim(c.Type, " ") != "" {
		updates["type"] = strings.Trim(c.Type, " ")
	}

	return ranking.Update(c.ID, updates)
}

func transfromSRK(mc ranking.Config) (srk.Config, error) {
	title := make(map[string]string)
	err := json.Unmarshal([]byte(mc.Title), &title)
	if err != nil {
		return srk.Config{}, err
	}
	problems := []map[string]any{}
	err = json.Unmarshal([]byte(mc.Problem), &problems)
	if err != nil {
		return srk.Config{}, err
	}
	members := []map[string]any{}
	err = json.Unmarshal([]byte(mc.Member), &members)
	if err != nil {
		return srk.Config{}, err
	}
	markers := []map[string]any{}
	err = json.Unmarshal([]byte(mc.Marker), &markers)
	if err != nil {
		return srk.Config{}, err
	}
	series := []map[string]any{}
	err = json.Unmarshal([]byte(mc.Series), &series)
	if err != nil {
		return srk.Config{}, err
	}
	sorter := map[string]any{}
	err = json.Unmarshal([]byte(mc.Sorter), &sorter)
	if err != nil {
		return srk.Config{}, err
	}
	contributors := []string{}
	_ = json.Unmarshal([]byte(mc.Contributor), &contributors)

	return srk.Config{
		ID:           mc.ID,
		UniqueKey:    mc.UniqueKey,
		Title:        title,
		StartAt:      mc.StartAt,
		Duration:     *srk.NewDuration(mc.EndAt.Sub(mc.StartAt)),
		Frozen:       *srk.NewDuration(time.Duration(mc.Frozen) * time.Millisecond),
		UnfrozenAt:   mc.UnfrozenAt,
		Problems:     problems,
		Members:      members,
		Markers:      markers,
		Series:       series,
		Sorter:       sorter,
		Contributors: contributors,
		Type:         mc.Type,
	}, nil
}

func srkTransfrom(ct srk.Config) (ranking.Config, error) {
	title, _ := json.Marshal(ct.Title)
	problems, _ := json.Marshal(ct.Problems)
	members, _ := json.Marshal(ct.Members)
	markers, _ := json.Marshal(ct.Markers)
	series, _ := json.Marshal(ct.Series)
	sorter, _ := json.Marshal(ct.Sorter)
	contributor, _ := json.Marshal(ct.Contributors)

	duration, err := ct.Duration.Duration()
	if err != nil {
		return ranking.Config{}, err
	}
	frozen, err := ct.Frozen.Duration()
	if err != nil {
		return ranking.Config{}, err
	}

	return ranking.Config{
		ID:          ct.ID,
		UniqueKey:   ct.UniqueKey,
		Title:       string(title),
		StartAt:     ct.StartAt,
		EndAt:       ct.StartAt.Add(duration),
		Frozen:      int64(frozen / time.Millisecond),
		UnfrozenAt:  ct.UnfrozenAt,
		Problem:     string(problems),
		Member:      string(members),
		Marker:      string(markers),
		Series:      string(series),
		Sorter:      string(sorter),
		Contributor: string(contributor),
		Type:        ct.Type,
	}, nil
}

func SetRecord(configID int64, records []srk.Record) error {
	rMap := make(map[string]map[string]string)
	for _, r := range records {
		k := fmt.Sprintf("%v:%v", configID, r.MemberID)
		if _, ok := rMap[k]; !ok {
			rMap[k] = make(map[string]string)
		}
		rMap[k][strconv.FormatInt(r.ID, 10)] = fmt.Sprintf("%v,%v,%v", r.ProblemID, r.Result, r.SubmissionTime)
	}

	ctx := context.Background()
	pipe := load.GetRedis().Pipeline()
	for k, v := range rMap {
		pipe.HMSet(ctx, k, v)
		pipe.Expire(ctx, k, 30*24*time.Hour)
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return err
	}

	go func() {
		SetRanking(configID)
		setRecords(configID, records)
	}()
	return nil
}

func GetRecord(contestID int64, memberIDs []string) (map[string][]srk.Record, error) {
	ctx := context.Background()
	pipe := load.GetRedis().Pipeline()
	for _, id := range memberIDs {
		pipe.HGetAll(ctx, fmt.Sprintf("%v:%v", contestID, id))
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	memberRecords := make(map[string][]srk.Record)
	for i, cmd := range cmds {
		rs, err := cmd.(*redis.MapStringStringCmd).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		memberID := memberIDs[i]
		records := make([]srk.Record, 0, len(rs))
		for k, v := range rs {
			id, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				return nil, err
			}
			v := strings.Split(v, ",")
			sulotion, err := strconv.ParseInt(v[2], 10, 64)
			if err != nil {
				return nil, err
			}
			records = append(records, srk.Record{
				ID:             id,
				MemberID:       memberID,
				ProblemID:      v[0],
				Result:         v[1],
				SubmissionTime: sulotion,
			})
		}
		memberRecords[memberID] = records
	}
	return memberRecords, nil
}

func SetRanking(id int64) {
	srk, err := GetSRKRank(id)
	if err != nil {
		return
	}
	if v, ok := Ranking[id]; ok {
		_ = v.Swap(srk)
		return
	}
	v := &atomic.Value{}
	v.Store(srk)
	Ranking[id] = v
}

func SetContestRecord(id int64, records []srk.Record) {
	// 需要 redis 发布
}

func GetRankingByConfigID(id int64) (string, error) {
	if _, ok := Ranking[id]; !ok {
		SetRanking(id)
	}

	r, ok := Ranking[id]
	if !ok {
		return "", errcode.NoResultErr
	}
	rankStr, ok := r.Load().(string)
	if !ok || rankStr == "" {
		return "", errcode.NoResultErr
	}

	return rankStr, nil
}

func GetRecordByConfigID(id int64, writer http.ResponseWriter, req *http.Request) error {
	wsConn, err := ws.WSHandler(writer, req, nil)
	if err != nil {
		return errcode.ServerErr
	}
	ws.NewRecordConn(id, wsConn)
	return nil
}

func GetSRKRank(contestID int64) (string, error) {
	ct, err := ranking.GetConfigByID(contestID)
	if err != nil {
		return "", err
	}
	if ct == nil {
		return "", errcode.NoResultErr
	}
	sc, err := transfromSRK(*ct)
	if err != nil {
		return "", err
	}

	memberIDs := make([]string, 0, len(sc.Members))
	for _, m := range sc.Members {
		memberIDs = append(memberIDs, m["id"].(string))
	}
	memberRecords, err := GetRecord(contestID, memberIDs)
	if err != nil {
		return "", err
	}

	srkRank := map[string]interface{}{
		"type":         sc.Type,
		"version":      getVersion(),
		"sorter":       sc.Sorter,
		"contributors": sc.Contributors,
		"series":       sc.Series,
		"contest":      getContest(sc.Title, sc.StartAt, sc.Duration, sc.Frozen),
		"problems":     sc.Problems,
		"rows":         getRows(sc, memberRecords),
		"markers":      sc.Markers,
		"_now":         time.Now().Format(time.RFC3339),
	}

	v, err := json.Marshal(srkRank)
	return string(v), err
}

func getVersion() string {
	return "0.3.0"
}

func getContest(title map[string]string, startAt time.Time, duration, frozen srk.Duration) map[string]interface{} {
	return map[string]interface{}{
		"title":          title,
		"startAt":        startAt.Format(time.RFC3339),
		"duration":       duration,
		"frozenDuration": frozen,
	}
}

type row struct {
	allTime  int64
	value    int64
	user     map[string]any
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

func getRows(sc srk.Config, memberRecords map[string][]srk.Record) []map[string]interface{} {
	cfg := sc.Sorter["config"]
	penalty := 20 * 60 // 默认罚时 20 分钟
	noPenalty := []string{SR_FirstBlood, SR_Accepted, SR_CompilationError, SR_UnknownError, SR_Frozen}
	if cfg, ok := cfg.(map[string]interface{}); ok {
		if p, ok := cfg["penalty"].([]interface{}); ok {
			d := &srk.Duration{p[0], p[1]}
			if t, err := d.Duration(); err == nil {
				penalty = int(t / time.Second)
			}
		}
		if np, ok := cfg["noPenaltyResults"].([]string); ok {
			noPenalty = np
		}
	}
	rows := make([]row, 0, len(sc.Members))
	for _, member := range sc.Members {
		records, ok := memberRecords[member["id"].(string)]
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
			return records[i].SubmissionTime < records[j].SubmissionTime
		})
		isSolutions := make(map[string]bool) // 存储题目是否已经被解决
		solutionMap := make(map[string][]solution)
		for _, r := range records {
			if isSolutions[r.ProblemID] {
				continue
			}
			d, _ := sc.Duration.Duration()
			f, _ := sc.Frozen.Duration()
			if d-f <= time.Duration(r.SubmissionTime)*time.Second {
				solutionMap[r.ProblemID] = append(solutionMap[r.ProblemID], solution{
					result: "?",
					time:   r.SubmissionTime,
				})
				continue
			}

			solutionMap[r.ProblemID] = append(solutionMap[r.ProblemID], solution{
				result: r.Result,
				time:   r.SubmissionTime,
			})
			if r.Result == SR_FirstBlood || r.Result == SR_Accepted {
				isSolutions[r.ProblemID] = true
			}
		}

		var allTime, value int64
		stats := make([]map[string]interface{}, len(sc.Problems))
		for i, p := range sc.Problems {
			solution, ok := solutionMap[p["alias"].(string)]
			if !ok {
				stats[i] = map[string]interface{}{
					"result": nil,
					"time":   []interface{}{0, "s"},
					"tries":  0,
				}
				continue
			}

			sols := make([]map[string]interface{}, 0, len(solutionMap[p["alias"].(string)]))
			for _, s := range solution {
				sols = append(sols, map[string]interface{}{
					"result": s.result,
					"time":   []interface{}{s.time, "s"},
				})
			}
			sLen := len(solution)
			sres := solution[sLen-1].result
			if sres != SR_FirstBlood && sres != SR_Accepted && sres != SR_Frozen {
				sres = SR_Rejected
			}
			stats[i] = map[string]interface{}{
				"result":    sres,
				"tries":     sLen,
				"solutions": sols,
			}
			if isSolutions[p["alias"].(string)] {
				stats[i]["time"] = []interface{}{solution[sLen-1].time, "s"}
				penaltyTries := 0
				for _, s := range solution {
					if !util.ContainsInSlice(noPenalty, s.result) {
						penaltyTries += 1
					}
				}
				pTime := solution[sLen-1].time + int64(penalty*penaltyTries)
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
			return rows[i].allTime < rows[j].allTime
		}
		return rows[i].value > rows[j].value
	})

	rs := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		user := row.user
		rs = append(rs, map[string]interface{}{
			"score":    map[string]interface{}{"value": row.value, "time": []interface{}{row.allTime, "s"}},
			"statuses": row.statuses,
			"user":     user,
		})
	}
	return rs
}

func setRecords(id int64, records []srk.Record) {
	isExist := make(map[string]bool)
	memberIDs := make([]string, 0, len(records))
	for _, r := range records {
		if isExist[r.MemberID] {
			continue
		}
		memberIDs = append(memberIDs, r.MemberID)
		isExist[r.MemberID] = true
	}
	memberRecords, err := GetRecord(id, memberIDs)
	if err != nil {
		return
	}

	ctx := context.Background()
	for _, record := range records {
		rs, ok := memberRecords[record.MemberID]
		if !ok {
			continue
		}
		solvedMap := make(map[string]bool)
		for _, r := range rs {
			if r.Result == SR_Accepted || r.Result == SR_FirstBlood {
				solvedMap[r.ProblemID] = true
			}
		}

		pubsub.Publish(ctx, fmt.Sprintf("ws:%v", id), ws.ScorllRecord{
			ID:        record.ID,
			ProblemID: record.ProblemID,
			MemberID:  record.MemberID,
			Result:    record.Result,
			Solved:    int8(len(solvedMap)),
		})
	}
}
