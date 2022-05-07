// Package collector loads exchange's orderbook and transactions datas.
package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/neosouler7/GObserver/kbt"
	"github.com/neosouler7/GObserver/upb"
)

var (
	ObMap *sync.Map
)

func collector() {
	// TODO. according to pairs, calculate & store on collector's ObMap
	for {
		fmt.Println(upb.ObMap)
		fmt.Println(kbt.ObMap)
		time.Sleep(time.Second)
	}
}

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
