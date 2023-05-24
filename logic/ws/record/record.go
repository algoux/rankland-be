package record

import (
	"context"
	"fmt"
	"rankland/load"
	"rankland/logic/srk"
	"rankland/model/ranking"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetResultAndSolved(id, rid int64, problemID, memberID string) (string, int8, bool) {
	memberRecords, err := GetRecord(id, memberID)
	if err != nil {
		return "", 0, false
	}
	config, err := ranking.GetConfigByID(id)
	if err != nil {
		return "", 0, false
	}
	t := config.EndAt.Sub(config.StartAt) - time.Duration(config.Frozen)*time.Millisecond

	sort.Slice(memberRecords, func(i, j int) bool {
		return memberRecords[i].ID < memberRecords[j].ID
	})
	result := ""
	solved := int8(0)
	problemSolved := make(map[string]bool)
	for _, r := range memberRecords {
		if problemSolved[r.ProblemID] {
			continue
		}

		if t <= time.Duration(r.SubmissionTime)*time.Second {
			if r.ID == rid {
				result = "?"
			}
			continue
		}

		if r.Result == "FB" || r.Result == "AC" {
			problemSolved[r.ProblemID] = true
			if t > time.Duration(r.SubmissionTime)*time.Second {
				solved += 1
			}
		}

		if r.ID == rid {
			result = r.Result
		}
	}
	if result == "" {
		return "", 0, false
	}

	return result, solved, true
}

func GetRecord(rankingID int64, memberID string) ([]srk.Record, error) {
	ctx := context.Background()
	memberRecords := make([]srk.Record, 0)
	cmd := load.GetRedis().HGetAll(ctx, fmt.Sprintf("%v:%v", rankingID, memberID))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	for k, v := range cmd.Val() {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}
		v := strings.Split(v, ",")
		sulotion, err := strconv.ParseInt(v[2], 10, 64)
		if err != nil {
			return nil, err
		}
		memberRecords = append(memberRecords, srk.Record{
			ID:             id,
			MemberID:       memberID,
			ProblemID:      v[0],
			Result:         v[1],
			SubmissionTime: sulotion,
		})
	}

	return memberRecords, nil
}
