package updater

import (
	"fmt"
	"time"
)

func Start() {
	fmt.Println("updater start")
	for {
		fmt.Println("updater")
		time.Sleep(time.Second * 2)
	}
}
