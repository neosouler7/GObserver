// Package upb provides upbit's ob & tx datas.
package upb

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
	wsConn *websocket.Conn
	ObMap  sync.Map
	TxMap  sync.Map
)

// Sends websocket ping message.
func ping() {
	for {
		ws.Ping(wsConn)
		// ws.SendMsg(wsConn, "PING")
		time.Sleep(time.Second * 5)
	}
}

// Subscribes upb's orderbook & transaction.
func subscribe() {
	time.Sleep(time.Second * 1)

	var streamSlice []string
	for _, pair := range utils.GetPairs(utils.UPB) {
		var pairInfo = strings.Split(pair, ":")
		market, symbol := strings.ToUpper(pairInfo[0]), strings.ToUpper(pairInfo[1])

		streamSlice = append(streamSlice, fmt.Sprintf("'%s-%s'", market, symbol))
	}
	uuid := uuid.NewString()
	streams := strings.Join(streamSlice, ",")
	msg := fmt.Sprintf("[{'ticket':'%s'}, {'type': 'trade', 'codes': [%s]}, {'type': 'orderbook', 'codes': [%s]}]", uuid, streams, streams)
	ws.SendMsg(wsConn, msg)
}

// Receives upb's orderbook & transaction.
func receive() {
	for {
		_, msgBytes, err := wsConn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(string(msgBytes), "status") {
			fmt.Println("PONG") // {"status":"UP"}
		} else {
			var rJson interface{}
			utils.Bytes2Json(msgBytes, &rJson)
			subject := rJson.(map[string]interface{})["type"].(string)

			switch subject {
			case "orderbook":
				fmt.Println("upb orderbook rcv")
				t := utils.GetTaker(utils.UPB, rJson.(map[string]interface{}))
				obKey := fmt.Sprintf("%s:%s:%s", t.Exchange, t.Market, t.Symbol)
				ObMap.Store(obKey, fmt.Sprintf("%s|%s", t.AskPrice, t.BidPrice))

			case "trade":
				// TODO.
				// map[ask_bid:BID change:FALL change_price:49000 code:KRW-ETH prev_closing_price:3.456e+06 sequential_id:1.651975281000001e+15
				// stream_type:REALTIME timestamp:1.651975281892e+12 trade_date:2022-05-08 trade_price:3.407e+06 trade_time:02:01:21
				// trade_timestamp:1.651975281e+12 trade_volume:0.29168568 type:trade]
				fmt.Println("upb transaction rcv")
			}
		}
	}
}

// Wakes up upb's websocket logic.
func Start() {
	log.Printf("collector - upb called.")
	wsConn = ws.GetConn(utils.UPB)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ping() // ping
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		subscribe() // subscribe websocket stream
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		receive() // receive websocket msg
	}()

	wg.Wait()
}
