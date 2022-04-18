package tg

import (
	"fmt"
	"log"
)

func Start() {
	fmt.Println("tg called")
}

func HandleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
