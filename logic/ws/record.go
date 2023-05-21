package ws

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"rankland/logic/pubsub"
	"sync"

	"github.com/gorilla/websocket"
)

var syncRecord map[int64]*RecordConn

type RecordConn struct {
	id    int64
	once  *sync.Once
	conns map[*websocket.Conn]bool
}

func setRecordConn(id int64, conn *websocket.Conn) {
	if _, ok := syncRecord[id]; !ok {
		syncRecord[id] = &RecordConn{
			id:    id,
			once:  &sync.Once{},
			conns: make(map[*websocket.Conn]bool),
		}
	}
	sr := syncRecord[id]
	sr.conns[conn] = true
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
		json.Unmarshal([]byte(str), r)
		buf := &bytes.Buffer{}
		// id, problemID, memberID, result, solved
		buf.Write([]byte{8, byte(len(r.ProblemID)), byte(len(r.MemberID)), byte(len(r.Result)), 1})
		binary.Write(buf, binary.BigEndian, r.ID)
		binary.Write(buf, binary.BigEndian, r.ProblemID)
		binary.Write(buf, binary.BigEndian, r.MemberID)
		binary.Write(buf, binary.BigEndian, r.Result)
		binary.Write(buf, binary.BigEndian, r.Solved)

		for conn := range sr.conns {
			conn.WriteMessage(1, buf.Bytes())
			// conn.WriteMessage(1, []byte(str))
		}
	}

	for conn := range sr.conns {
		conn.Close()
	}
	defer delete(syncRecord, id)
}

func NewRecordConn(id int64, conn *websocket.Conn) {
	if syncRecord == nil {
		syncRecord = make(map[int64]*RecordConn)
	}
	setRecordConn(id, conn)
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
