package cpu

import (
	"fmt"
	"time"
)

func processor() {
	fmt.Println("processor start")
	for {
		fmt.Println("processor")
		time.Sleep(time.Second * 3)
	}
}
