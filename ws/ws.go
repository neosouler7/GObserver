// Package ws provides websocket functions.
package ws

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	w *websocket.Conn
)

type hostPath struct {
	host string
	path string
}

// Send websocket payload message.
func SendMsg(conn *websocket.Conn, msg string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Fatalln(err)
	}
}

// Sends websocket ping message for alive connection.
func Ping(conn *websocket.Conn) {
	err := conn.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		log.Fatalln(err)
	}
}

// Returns host & path info by exchange.
func (h *hostPath) getHostPath(exchange string) {
	switch exchange {
	case "kbt":
		h.host = "ws.korbit.co.kr"
		h.path = "/v1/user/push"
	case "upb":
		h.host = "api.upbit.com"
		h.path = "/websocket/v1"
	}
}

// Returns websocket connection by exchange.
func GetConn(exchange string) *websocket.Conn {
	h := &hostPath{}
	h.getHostPath(exchange)

	u := url.URL{Scheme: "wss", Host: h.host, Path: h.path}
	wPointer, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	w = wPointer
	return w
}
