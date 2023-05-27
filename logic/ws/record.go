package ws

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"rankland/load"
	"rankland/logic/pubsub"
	"rankland/logic/srk"
	"rankland/model/ranking"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var syncRecord map[int64]*RecordConn

type RecordConn struct {
	id    int64
	once  *sync.Once
	conns map[*websocket.Conn]bool
}

func setRecordConn(id int64, conn *websocket.Conn, isAdmin bool) {
	if _, ok := syncRecord[id]; !ok {
		syncRecord[id] = &RecordConn{
			id:    id,
			once:  &sync.Once{},
			conns: make(map[*websocket.Conn]bool),
		}
	}
	sr := syncRecord[id]
	sr.conns[conn] = isAdmin
	sr.once.Do(func() {
		go writeRecord(id)
	})
}

type ScorllRecord struct {
	ID        int64
	ProblemID string
	MemberID  string
	Result    string
	Solved    int8
}

func (sr ScorllRecord) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(sr)
	return bytes, err
}

func writeRecord(id int64) {
	ctx := context.Background()
	channel := make(chan string, 100)
	go pubsub.Subscribe(ctx, fmt.Sprintf("ws:%v", id), channel)

	sr := syncRecord[id]
	for str := range channel {
		r := &ScorllRecord{}
		err := json.Unmarshal([]byte(str), r)
		if err != nil {
			continue
		}

		unfronzenRet, ok := GetResultAndSolved(id, r.ID, r.ProblemID, r.MemberID, true)
		if !ok {
			continue
		}
		fronzenRet, ok := GetResultAndSolved(id, r.ID, r.ProblemID, r.MemberID, false)
		if !ok {
			continue
		}

		for conn, isUnfronzen := range sr.conns {
			if isUnfronzen {
				conn.WriteMessage(websocket.BinaryMessage, unfronzenRet)
			} else {
				conn.WriteMessage(websocket.BinaryMessage, fronzenRet)
			}
		}
	}

	for conn := range sr.conns {
		conn.Close()
	}
	defer delete(syncRecord, id)
}

func NewRecordConn(id int64, conn *websocket.Conn, isAdmin bool) {
	if syncRecord == nil {
		syncRecord = make(map[int64]*RecordConn)
	}
	setRecordConn(id, conn, isAdmin)
}

// type RecordCli struct {
// 	id   int64
// 	conn *websocket.Conn
// }

// func NewRecordCli(id int64, conn *websocket.Conn) *RecordCli {
// 	return &RecordCli{
// 		id:   id,
// 		conn: conn,
// 	}
// }

// func (r *RecordCli) Read() (string, error) {
// 	_, d, err := r.conn.ReadMessage()
// 	return string(d), err
// }

// func (r *RecordCli) Write(d []byte) error {
// 	return r.conn.WriteMessage(1, d)
// }

// func (r *RecordCli) Close() error {
// 	return r.conn.Close()
// }

func GetResultAndSolved(id, rid int64, problemID, memberID string, isUnfrozen bool) ([]byte, bool) {
	memberRecords, err := GetRecord(id, memberID)
	if err != nil {
		return nil, false
	}
	config, err := ranking.GetConfigByID(id)
	if err != nil {
		return nil, false
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

		if t <= time.Duration(r.SubmissionTime)*time.Second && !isUnfrozen {
			if r.ID == rid {
				result = "?"
				continue
			}
		}

		if r.Result == "FB" || r.Result == "AC" {
			problemSolved[r.ProblemID] = true
			if t > time.Duration(r.SubmissionTime)*time.Second || isUnfrozen {
				solved += 1
			}
		}

		if r.ID == rid {
			result = r.Result
		}
	}
	if result == "" {
		return nil, false
	}

	buf := &bytes.Buffer{}
	// id, problemID, memberID, result, solved
	buf.Write([]byte{5, 8, byte(len(problemID)), byte(len(memberID)), byte(len(result)), 1})
	binary.Write(buf, binary.BigEndian, rid)
	buf.Write([]byte(problemID))
	buf.Write([]byte(memberID))
	buf.Write([]byte(result))
	binary.Write(buf, binary.BigEndian, solved)

	return buf.Bytes(), true
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
