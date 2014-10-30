package main

import (
	"fmt"

	"github.com/abourget/arigo"
)

func main() {
	a := arigo.NewARI("asterisk", "asterisk", "localhost", 8088)
	err := a.Connect("hello-world")
	fmt.Println("Connecting... ", err)

	go a.HandleReceive()

	for {
		select {
		case msg := <- a.ReceiveChan:
			fmt.Printf("Received a message: %#v\n", msg)
		}
	}
}
