// Package collector loads exchange's ob & tx and stores on syncMap.
package collector

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/neosouler7/GObserver/config"
	"github.com/neosouler7/GObserver/kbt"
	"github.com/neosouler7/GObserver/upb"
	"github.com/neosouler7/GObserver/utils"
)

var (
	ObMap   sync.Map
	obTsMap sync.Map
	TxMap   sync.Map
	txTsMap sync.Map
)

// Collects exchange's data and saves to collector's obMap.
func collector() {
	// TODO. according to pairs, calculate & store on collector's ObMap & TxMap.
	var combs []string
	for _, exchange := range config.GetExchanges() {
		for _, pair := range config.GetPairs(exchange) {
			combs = append(combs, fmt.Sprintf("%s:%s", exchange, pair))
		}
	}

	for {
		for _, comb := range combs {
			var ob, prevObTsStr interface{}
			var ok bool
			exchange := strings.Split(comb, ":")[0]

			switch exchange {
			case config.UPB:
				ob, ok = upb.ObMap.Load(comb)
			case config.KBT:
				ob, ok = kbt.ObMap.Load(comb)
			}

			if !ok {
				log.Printf("waiting collector's obMap to be stored.")
				time.Sleep(time.Second)
				continue
			}

			tsStr := ob.(*utils.Taker).Timestamp
			ts, _ := strconv.ParseInt(tsStr, 10, 64)
			prevObTsStr, ok = obTsMap.Load(comb)
			if !ok {
				log.Printf("init store on obTsMap of comb")
				obTsMap.Store(comb, tsStr) // init store ob's timestamp.
				continue
			}

			prevObTs, _ := strconv.ParseInt(prevObTsStr.(string), 10, 64)
			if ts-prevObTs > 0 { // store when only new ts is bigger than previous.
				ObMap.Store(comb, ob)
				log.Printf(fmt.Sprintf("OB %s collector stored", comb))
				time.Sleep(time.Millisecond * 100)
			}
		}
	}
}

// Starts collector logic of "cpu".
func Start() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		collector()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		upb.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		kbt.Start()
	}()

	wg.Wait()
}
