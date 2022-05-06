// Package main runs cpu / db / tg.
package main

import (
	"sync"

	"github.com/neosouler7/GObserver/db"
	"github.com/neosouler7/GObserver/tg"
)

// Runs packages as goroutine.
func main() {
	var wg sync.WaitGroup
	wg.Add(2) // cpu, db, tg

	// // go func() {
	// // 	defer wg.Done()
	// // 	cpu.Collect()
	// // }()

	go func() {
		defer wg.Done()
		tg.Start()
	}()

	go func() {
		defer wg.Done()
		db.Start()
	}()

	wg.Wait()
}
