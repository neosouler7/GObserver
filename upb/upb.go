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
	obConn *websocket.Conn
	txConn *websocket.Conn
	ObMap  sync.Map
	TxMap  sync.Map
)

const (
	pingMsg string = "PING"
)

func obPing() {
	for {
		ws.SendMsg(obConn, pingMsg)
		time.Sleep(time.Second * 5)
	}
}

func obSub(pairs []string) {
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

	ws.SendMsg(obConn, msg)
}

func obRcv() {
	for {
		_, msgBytes, err := obConn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(string(msgBytes), "status") {
			fmt.Println("PONG") // {"status":"UP"}
		} else {
			var rJson interface{}
			utils.Bytes2Json(msgBytes, &rJson)

			t := utils.GetTaker(utils.UPB, rJson.(map[string]interface{}))
			obKey := fmt.Sprintf("%s:%s:%s", t.Exchange, t.Market, t.Symbol)
			ObMap.Store(obKey, fmt.Sprintf("%s|%s", t.AskPrice, t.BidPrice))
		}
	}
}

// Subscribes upb's orderbook & transaction.
func Start() {
	log.Printf("collector - upb called.")
	obConn = ws.GetConn(utils.UPB, utils.OB)
	txConn = ws.GetConn(utils.UPB, utils.TX)

	var wg sync.WaitGroup

	// orderbook
	wg.Add(1)
	go obPing() // ping

	wg.Add(1)
	pairs := []string{"krw:btc", "krw:eth", "krw:xrp"}
	go func() {
		obSub(pairs) // subscribe websocket stream
		wg.Done()
	}()

	wg.Add(1)
	go obRcv() // receive websocket msg

	// TODO. transaction

	wg.Wait()
}
