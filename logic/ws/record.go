package ws

import (
	"context"
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

func writeRecord(id int64) {
	ctx := context.Background()
	channel := make(chan []byte, 100)
	pubsub.Subscribe(ctx, fmt.Sprintf("ws:%v", id), channel)

	sr := syncRecord[id]
	for b := range channel {
		for conn := range sr.conns {
			conn.WriteMessage(1, b)
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
