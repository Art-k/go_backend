package src

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	Ws *websocket.Conn

	err_ws   error
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connected to socket")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	Ws, err_ws = upgrader.Upgrade(w, r, nil)
	if err_ws != nil {
		if _, ok := err_ws.(websocket.HandshakeError); !ok {
			log.Println(err_ws)
		}
		return
	}
}
