package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

func main() {
	a := ari.NewARI("asterisk", "asterisk", "localhost", 8088, "hello-world")
	rest := a.GetREST()

	receiveChan := a.LaunchListener()

	for {
		select {
		case msg := <-receiveChan:
			switch m := msg.(type) {
			case *ari.AriConnected:
				fmt.Println("Ok, connected, sending AsteriskInfo request")

				infos, err := rest.AsteriskInfoGet()
				if err != nil {
					fmt.Println("Couldn'get infos", err)
				} else {
					pretty.Printf("AsteriskInfos: %# v\n", infos)
				}

				if res, err := rest.AsteriskVariableGet("TRUNK"); err == nil {
					fmt.Println("Got variable:", res)
				}
			case ari.ChannelHangupRequest:
				fmt.Printf("Hangup for channel %s\n", m.Channel)
			default:
				pretty.Printf("Received some message: %# v\n", msg)
			}
		}

	}
}
