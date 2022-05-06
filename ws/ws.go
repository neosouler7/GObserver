package ws

import (
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	once sync.Once
	w    *websocket.Conn
)

type hostPath struct {
	host string
	path string
}

func SendMsg(conn *websocket.Conn, msg string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Fatalln(err)
	}
}

func GetConn(exchange, subject string) *websocket.Conn {
	// once.Do(func() {
	h := &hostPath{}
	h.getHostPath(exchange, subject)

	u := url.URL{Scheme: "wss", Host: h.host, Path: h.path}
	wPointer, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	w = wPointer
	// })
	return w
}

func (h *hostPath) getHostPath(exchange, subject string) {
	switch exchange {
	case "kbt":
		h.host = "ws.korbit.co.kr"
		switch subject {
		case "orderbook":
			h.path = "/v1/user/push"
		case "transaction":
			h.path = "/v1/user/push" // TODO.
		}
	case "upb":
		h.host = "api.upbit.com"
		switch subject {
		case "orderbook":
			h.path = "/websocket/v1"
		case "transaction":
			h.path = "/websocket/v1" // TODO.
		}
	}
}
