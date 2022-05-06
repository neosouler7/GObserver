package cpu

import (
	"fmt"
	"time"
)

func updater() {
	fmt.Println("updater start")
	for {
		fmt.Println("updater")
		time.Sleep(time.Second * 2)
	}
}
