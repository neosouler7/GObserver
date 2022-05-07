package kbt

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/neosouler7/GObserver/utils"
	"github.com/neosouler7/GObserver/ws"
)

var (
	obConn *websocket.Conn
	txConn *websocket.Conn
	ObMap  sync.Map
	TxMap  sync.Map
)

func kbtObSub(pairs []string) {
	time.Sleep(time.Second * 1)
	var streamSlice []string
	for _, pair := range pairs {
		var pairInfo = strings.Split(pair, ":")
		market, symbol := strings.ToLower(pairInfo[0]), strings.ToLower(pairInfo[1])

		streamSlice = append(streamSlice, fmt.Sprintf("%s_%s", symbol, market))
	}

	ts := time.Now().UnixNano() / 100000 / 10
	streams := fmt.Sprintf("\"orderbook:%s\"", strings.Join(streamSlice, ","))
	msg := fmt.Sprintf("{\"accessToken\": \"null\", \"timestamp\": \"%d\", \"event\": \"korbit:subscribe\", \"data\": {\"channels\": [%s]}}", ts, streams)

	ws.SendMsg(obConn, msg)
}

func obRcv() {
	pairs := []string{"krw:btc", "krw:eth", "krw:xrp"}
	for {
		_, msgBytes, err := obConn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(string(msgBytes), "connected") {
			kbtObSub(pairs) // just once
		} else if strings.Contains(string(msgBytes), "subscribe") {
			continue
		} else if strings.Contains(string(msgBytes), "push-orderbook") {
			var rJson interface{}
			utils.Bytes2Json(msgBytes, &rJson)

			t := utils.GetTaker(utils.KBT, rJson.(map[string]interface{}))
			obKey := fmt.Sprintf("%s:%s:%s", t.Exchange, t.Market, t.Symbol)
			ObMap.Store(obKey, fmt.Sprintf("%s|%s", t.AskPrice, t.BidPrice))
		} else {
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// Subscribes kbt's orderbook & transaction.
func Start() {
	log.Printf("collector - kbt called.")
	obConn = ws.GetConn(utils.KBT, utils.OB)
	txConn = ws.GetConn(utils.KBT, utils.TX)

	var wg sync.WaitGroup

	// orderbook
	wg.Add(1)
	go obRcv() // receive websocket msg

	// TODO. transaction

	wg.Wait()
}
