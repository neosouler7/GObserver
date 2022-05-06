// Package cpu collects & processes & updates informations from cryto-exchange markets to database.
package cpu

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/neosouler7/GObserver/utils"
)

const (
	UPB string = "upb"
	KBT string = "kbt"
	OB  string = "orderbook"
	TX  string = "transaction"
)

var (
	// targetVolumeMap = make(map[string]string)
	ObMap *sync.Map
)

type taker struct {
	exchange string
	market   string
	symbol   string
	askPrice string
	bidPrice string
}

// func getTargetVolume(exchange, market, symbol string) string {
// 	return targetVolumeMap[fmt.Sprintf("%s:%s:%s", exchange, market, symbol)]
// }

func getTargetPrice(volume string, orderbook interface{}) string {
	currentVolume := 0.0
	targetVolume, err := strconv.ParseFloat(volume, 64)
	if err != nil {
		log.Fatalln(err)
	}

	obSlice := orderbook.([]interface{})
	for _, ob := range obSlice {
		obInfo := ob.([2]string)
		volume, err := strconv.ParseFloat(obInfo[1], 64)
		if err != nil {
			log.Fatalln(err)
		}

		currentVolume += volume
		if currentVolume >= targetVolume {
			return obInfo[0]
		}
	}
	return obSlice[len(obSlice)-1].([2]string)[0]
}

func getTaker(exchange string, rJson map[string]interface{}) *taker {
	t := &taker{}
	var market, symbol string
	var askSlice, bidSlice []interface{}

	switch exchange {
	case UPB:
		s := strings.Split(rJson["code"].(string), "-")
		market, symbol = strings.ToLower(s[0]), strings.ToLower(s[1])

		obs := rJson["orderbook_units"].([]interface{})

		for _, ob := range obs {
			o := ob.(map[string]interface{})
			ask := [2]string{fmt.Sprintf("%f", o["ask_price"]), fmt.Sprintf("%f", o["ask_size"])}
			bid := [2]string{fmt.Sprintf("%f", o["bid_price"]), fmt.Sprintf("%f", o["bid_size"])}
			askSlice = append(askSlice, ask)
			bidSlice = append(bidSlice, bid)
		}

	case KBT:
		s := strings.Split(rJson["data"].(map[string]interface{})["currency_pair"].(string), "_")
		market, symbol = s[1], s[0]

		rData := rJson["data"]
		askResponse := rData.(map[string]interface{})["asks"].([]interface{})
		bidResponse := rData.(map[string]interface{})["bids"].([]interface{})

		for i := 0; i < utils.Min(len(askResponse), len(bidResponse))-1; i++ {
			askR, bidR := askResponse[i].(map[string]interface{}), bidResponse[i].(map[string]interface{})
			ask := [2]string{askR["price"].(string), askR["amount"].(string)}
			bid := [2]string{bidR["price"].(string), bidR["amount"].(string)}
			askSlice = append(askSlice, ask)
			bidSlice = append(bidSlice, bid)
		}
	}
	targetVolume := utils.GetVolumeMap(exchange)[fmt.Sprintf("%s:%s", market, symbol)]

	t.exchange = exchange
	t.market = market
	t.symbol = symbol
	t.askPrice = getTargetPrice(targetVolume, askSlice)
	t.bidPrice = getTargetPrice(targetVolume, bidSlice)

	return t
}

// Runs packages as goroutine.
func Start() {
	// // TODO
	// targetVolumeMap["upb:krw:btc"] = "0.1"
	// targetVolumeMap["upb:krw:eth"] = "1.2"
	// targetVolumeMap["upb:krw:xrp"] = "6250"
	// targetVolumeMap["kbt:krw:btc"] = "0.1"
	// targetVolumeMap["kbt:krw:eth"] = "1.2"
	// targetVolumeMap["kbt:krw:xrp"] = "6250"

	var temp sync.Map
	ObMap = &temp

	var wg sync.WaitGroup

	wg.Add(2) // processor, updater
	go func() {
		defer wg.Done()
		processor()
	}()

	go func() {
		defer wg.Done()
		updater()
	}()

	// collector
	wg.Add(1)
	go func() {
		defer wg.Done()
		upb()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		kbt()
	}()

	wg.Wait()
}
