// Package kbt provides korbit's ob & tx datas.
package kbt

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/neosouler7/GObserver/config"
	"github.com/neosouler7/GObserver/utils"
	"github.com/neosouler7/GObserver/ws"
)

var (
	wsConn *websocket.Conn
	ObMap  sync.Map
	TxMap  sync.Map
)

// Subscribes kbt's orderbook & transaction.
func subscribe(pairs []string) {
	time.Sleep(time.Second * 1)

	var streamSlice []string
	for _, pair := range pairs {
		var pairInfo = strings.Split(pair, ":")
		market, symbol := strings.ToLower(pairInfo[0]), strings.ToLower(pairInfo[1])

		streamSlice = append(streamSlice, fmt.Sprintf("%s_%s", symbol, market))
	}

	ts := time.Now().UnixNano() / 100000 / 10
	streamsOb := fmt.Sprintf("\"orderbook:%s\"", strings.Join(streamSlice, ","))
	streamsTx := fmt.Sprintf("\"transaction:%s\"", strings.Join(streamSlice, ","))
	msg := fmt.Sprintf("{\"accessToken\": \"null\", \"timestamp\": \"%d\", \"event\": \"korbit:subscribe\", \"data\": {\"channels\": [%s,%s]}}", ts, streamsOb, streamsTx)

	ws.SendMsg(wsConn, msg)
}

// Receives kbt's orderbook & transaction.
func receive() {
	pairs := config.GetPairs(config.KBT)
	for {
		_, msgBytes, err := wsConn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(string(msgBytes), "connected") {
			subscribe(pairs) // just once
		} else if strings.Contains(string(msgBytes), "subscribe") {
			continue
		} else {
			var rJson interface{}
			utils.Bytes2Json(msgBytes, &rJson)
			subject := rJson.(map[string]interface{})["data"].(map[string]interface{})["channel"]

			switch subject {
			case "orderbook":
				log.Printf("kbt orderbook rcv.")
				t := utils.GetTaker(config.KBT, rJson.(map[string]interface{}))
				obKey := fmt.Sprintf("%s:%s:%s", t.Exchange, t.Market, t.Symbol)
				ObMap.Store(obKey, t)

			case "transaction":
				// TODO
				// map[data:map[amount:66.930000 channel:transaction currency_pair:xrp_krw price:754.1 taker:buy timestamp:1.651976579592e+12]
				// event:korbit:push-transaction timestamp:1.65197657967e+12]
				log.Printf("kbt transaction rcv.")
			}
		}

		// {
		// 	"accessToken": null,
		// 	"event": "korbit:push-transaction",
		// 	"timestamp" : 1389678052000,
		// 	"data":
		// 	{
		// 	  "channel": "transaction",
		// 	  "currency_pair": "btc_krw",
		// 	  "timestamp" : 1389678052000,
		// 	  "price" : "569000.7654835",
		// 	  "amount" : "0.01000001",
		// 	  "taker" : "buy"
		// 	}
		//   }
		// } else {
		// 	if err != nil {
		// 		log.Fatalln(err)
		// 	}
		// }
	}
}

// Wakes up kbt's websocket logic.
func Start() {
	log.Printf("collector - kbt called.")
	wsConn = ws.GetConn(config.KBT)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		receive() // receive websocket msg
	}()

	wg.Wait()
}
