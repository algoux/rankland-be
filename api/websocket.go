package api

import (
	"fmt"
	"net/http"
	"rankland/errcode"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 允许跨域
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Errors = append(c.Errors, errcode.ServerErr)
		return
	}
	go sendMsg(wsConn)
}

func sendMsg(ws *websocket.Conn) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		data := fmt.Sprintf("hello word!, index: %v, time: %v", i, time.Now().Format(time.RFC3339))
		ws.WriteMessage(1, []byte(data))
	}
	defer func(ws *websocket.Conn) {
		ws.Close()
	}(ws)
}
