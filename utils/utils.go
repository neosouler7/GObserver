// Package utils provides frequent used functions among project.
package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/neosouler7/GObserver/config"
)

// Returns if a is contained in s.
func Contains[T comparable](s []T, a T) bool {
	for _, b := range s {
		if b == a {
			return true
		}
	}
	return false
}

// Returns smaller value.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Converts any type as bytes.
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	err := encoder.Encode(i)
	if err != nil {
		log.Fatalln(err)
	}
	return aBuffer.Bytes()
}

// Converts bytes to json.
func Bytes2Json(data []byte, i interface{}) {
	r := bytes.NewReader(data)
	err := json.NewDecoder(r).Decode(i)
	if err != nil {
		log.Fatalln(err)
	}
}

type taker struct {
	Exchange string
	Market   string
	Symbol   string
	AskPrice string
	BidPrice string
}

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

func GetTaker(exchange string, rJson map[string]interface{}) *taker {
	t := &taker{}
	var market, symbol string
	var askSlice, bidSlice []interface{}

	switch exchange {
	case config.UPB:
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

	case config.KBT:
		s := strings.Split(rJson["data"].(map[string]interface{})["currency_pair"].(string), "_")
		market, symbol = s[1], s[0]

		rData := rJson["data"]
		askResponse := rData.(map[string]interface{})["asks"].([]interface{})
		bidResponse := rData.(map[string]interface{})["bids"].([]interface{})

		for i := 0; i < Min(len(askResponse), len(bidResponse))-1; i++ {
			askR, bidR := askResponse[i].(map[string]interface{}), bidResponse[i].(map[string]interface{})
			ask := [2]string{askR["price"].(string), askR["amount"].(string)}
			bid := [2]string{bidR["price"].(string), bidR["amount"].(string)}
			askSlice = append(askSlice, ask)
			bidSlice = append(bidSlice, bid)
		}
	}
	targetVolume := config.GetVolumeMap(exchange)[fmt.Sprintf("%s:%s", market, symbol)]

	t.Exchange = exchange
	t.Market = market
	t.Symbol = symbol
	t.AskPrice = getTargetPrice(targetVolume, askSlice)
	t.BidPrice = getTargetPrice(targetVolume, bidSlice)

	return t
}
