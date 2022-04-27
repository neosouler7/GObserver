// Package main runs cpu / db / tg.
package main

import (
	"strconv"
	"time"

	"github.com/neosouler7/GObserver/alog"
)

func main() {
	go alog.Start()

	go func() {
		for i := 1; i < 20; i++ {
			n := strconv.Itoa(i)
			println(n)
			alog.LogC <- n
		}
	}()

	time.Sleep(1 * time.Second)
	// close(alog.LogC)

	// var wg sync.WaitGroup
	// wg.Add(1) // cpu, db, tg

	// // go func() {
	// // 	defer wg.Done()
	// // 	cpu.Collect()
	// // }()

	// // go func() {
	// // 	defer wg.Done()
	// // 	tg.Start()
	// // }()

	// go func() {
	// 	defer wg.Done()
	// 	db.Start()
	// }()

	// wg.Wait()
}
