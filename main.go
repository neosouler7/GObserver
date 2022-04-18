package main

import (
	"sync"

	"github.com/neosouler7/GObserver/cpu"
	"github.com/neosouler7/GObserver/db"
	"github.com/neosouler7/GObserver/tg"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3) // cpu, db, tg

	go func() {
		defer wg.Done()
		cpu.Collect()
	}()

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
