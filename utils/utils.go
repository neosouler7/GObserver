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

// Format timestamp as millisecond unit.
func formatTs(ts string) string {
	if len(ts) < 13 {
		add := strings.Repeat("0", 13-len(ts))
		return fmt.Sprintf("%s%s", ts, add)
	} else if len(ts) == 13 { // if millisecond
		return ts
	} else {
		return ts[:13]
	}
}

type Taker struct {
	Exchange  string
	Market    string
	Symbol    string
	AskPrice  string
	AskVolume string
	BidPrice  string
	BidVolume string
	Timestamp string
}

func getTargetPrice(volume string, orderbook interface{}) (string, string) {
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
			return obInfo[0], fmt.Sprintf("%f", currentVolume)
		}
	}
	return obSlice[len(obSlice)-1].([2]string)[0], fmt.Sprintf("%f", currentVolume)
}

func GetTaker(exchange string, rJson map[string]interface{}) *Taker {
	var market, symbol string
	var askSlice, bidSlice []interface{}
	var askPrice, askVolume, bidPrice, bidVolume, timestamp string

	switch exchange {
	case config.UPB:
		s := strings.Split(rJson["code"].(string), "-")
		market, symbol = strings.ToLower(s[0]), strings.ToLower(s[1])
		timestamp = fmt.Sprintf("%d", int(rJson["timestamp"].(float64)))

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
		timestamp = fmt.Sprintf("%d", int(rJson["data"].(map[string]interface{})["timestamp"].(float64)))

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

	askPrice, askVolume = getTargetPrice(targetVolume, askSlice)
	bidPrice, bidVolume = getTargetPrice(targetVolume, bidSlice)

	t := Taker{
		Exchange:  exchange,
		Market:    market,
		Symbol:    symbol,
		AskPrice:  askPrice,
		AskVolume: askVolume,
		BidPrice:  bidPrice,
		BidVolume: bidVolume,
		Timestamp: formatTs(timestamp),
	}
	return &t
}
