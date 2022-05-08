// Package updater loads processor's syncMap and stores to boltdb.
package updater

import (
	"fmt"
	"time"
)

// Starts updater logic of "cpu".
func Start() {
	fmt.Println("updater start")
	for {
		fmt.Println("updater")
		time.Sleep(time.Second * 2)
	}
}
