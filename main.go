// Package main runs cpu / db / tg.
package main

import (
	"runtime"
	"sync"

	"github.com/neosouler7/GObserver/collector"
	"github.com/neosouler7/GObserver/db"
	"github.com/neosouler7/GObserver/processor"
	"github.com/neosouler7/GObserver/tg"
	"github.com/neosouler7/GObserver/updater"
)

// Runs packages as goroutine.
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use max goroutines.

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		collector.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		processor.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		updater.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tg.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		db.Start()
	}()

	wg.Wait()
}
