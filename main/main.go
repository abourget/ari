package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

func main() {
	a := ari.NewARI("asterisk", "asterisk", "localhost", 8088)
	err := a.Connect("hello-world")
	fmt.Println("Connecting... ", err)

	go a.HandleReceive()

	for {
		select {
		case msg := <-a.ReceiveChan:
			pretty.Printf("Received a message: %# v\n", msg)
		}
	}
}
