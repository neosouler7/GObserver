// Package collector loads exchange's ob & tx and stores on syncMap.
package collector

import (
	"sync"
	"time"

	"github.com/neosouler7/GObserver/kbt"
	"github.com/neosouler7/GObserver/upb"
)

var (
	ObMap *sync.Map
	TxMap *sync.Map
)

// Collects exchange's data and saves to collector's obMap.
func collector() {
	// TODO. according to pairs, calculate & store on collector's ObMap
	for {
		// fmt.Println(upb.ObMap)
		// fmt.Println(kbt.ObMap)
		time.Sleep(time.Second)
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
