package processor

import (
	"fmt"
	"time"
)

func Start() {
	fmt.Println("processor start")
	for {
		fmt.Println("processor")
		time.Sleep(time.Second * 3)
	}
}
