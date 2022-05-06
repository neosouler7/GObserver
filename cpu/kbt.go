package cpu

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
	kbtObConn *websocket.Conn
	kbtTxConn *websocket.Conn
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

	ws.SendMsg(kbtObConn, msg)
}

func kbtObRcv() {
	pairs := []string{"krw:btc", "krw:eth", "krw:xrp"}
	for {
		_, msgBytes, err := kbtObConn.ReadMessage()
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

			t := getTaker(KBT, rJson.(map[string]interface{}))
			fmt.Println(t)
			obKey := fmt.Sprintf("%s:%s:%s", t.exchange, t.market, t.symbol)
			ObMap.Store(obKey, fmt.Sprintf("%s|%s", t.askPrice, t.bidPrice))
		} else {
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

// Subscribes kbt's orderbook & transaction.
func kbt() {
	log.Printf("collector-kbt called.")
	kbtObConn = ws.GetConn(KBT, OB)
	kbtTxConn = ws.GetConn(KBT, TX)

	var wg sync.WaitGroup

	// orderbook
	wg.Add(1)
	go kbtObRcv() // receive websocket msg

	// TODO. transaction
	wg.Wait()
}
