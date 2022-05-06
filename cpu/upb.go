package cpu

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/neosouler7/GObserver/utils"
	"github.com/neosouler7/GObserver/ws"
)

var (
	upbObConn *websocket.Conn
	upbTxConn *websocket.Conn
)

const (
	upbObPingMsg string = "PING"
)

func upbObPing() {
	for {
		ws.SendMsg(upbObConn, upbObPingMsg)
		time.Sleep(time.Second * 5)
	}
}

func upbObSub(pairs []string) {
	time.Sleep(time.Second)
	var streamSlice []string
	for _, pair := range pairs {
		var pairInfo = strings.Split(pair, ":")
		market, symbol := strings.ToUpper(pairInfo[0]), strings.ToUpper(pairInfo[1])

		streamSlice = append(streamSlice, fmt.Sprintf("'%s-%s'", market, symbol))
	}
	uuid := uuid.NewString()
	streams := strings.Join(streamSlice, ",")
	msg := fmt.Sprintf("[{'ticket':'%s'}, {'type': 'orderbook', 'codes': [%s]}]", uuid, streams)

	ws.SendMsg(upbObConn, msg)
}

func upbObRcv() {
	for {
		_, msgBytes, err := upbObConn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(string(msgBytes), "status") {
			fmt.Println("PONG") // {"status":"UP"}
		} else {
			var rJson interface{}
			utils.Bytes2Json(msgBytes, &rJson)

			t := getTaker(UPB, rJson.(map[string]interface{}))
			fmt.Println(t)
			obKey := fmt.Sprintf("%s:%s:%s", t.exchange, t.market, t.symbol)
			ObMap.Store(obKey, fmt.Sprintf("%s|%s", t.askPrice, t.bidPrice))
		}
	}
}

// Subscribes upb's orderbook & transaction.
func upb() {
	log.Printf("collector-upb called.")
	upbObConn = ws.GetConn(UPB, OB)
	upbTxConn = ws.GetConn(UPB, TX)

	var wg sync.WaitGroup

	// orderbook
	wg.Add(1)
	go upbObPing() // ping

	wg.Add(1)
	pairs := []string{"krw:btc", "krw:eth", "krw:xrp"}
	go func() {
		upbObSub(pairs) // subscribe websocket stream
		wg.Done()
	}()

	wg.Add(1)
	go upbObRcv() // receive websocket msg

	// TODO. transaction

	wg.Wait()
}
