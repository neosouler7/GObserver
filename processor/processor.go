// Package processor loads collector's datas and calculates hit count on syncMap.
package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/neosouler7/GObserver/collector"
)

var (
	HitCountMap *sync.Map
)

// Starts processor logic of "cpu".
func Start() {
	fmt.Println("processor start")
	for {
		// TODO. calculate hit logic
		fmt.Println("processor")
		fmt.Println(collector.ObMap)
		fmt.Println(collector.TxMap)
		time.Sleep(time.Second * 3)
	}
}
