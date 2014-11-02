package main

// Package main is a demo of feature and a sample program that uses `ari`, `rest` and `models`.

import (
	"fmt"

	"github.com/abourget/ari"
	ast "github.com/abourget/ari/models"
	"github.com/abourget/ari/rest"
	"github.com/kr/pretty"
)

func main() {
	a := ari.NewARI("asterisk", "asterisk", "localhost", 8088, "hello-world")
	a.Debug = true
	r := a.GetREST()
	r.Debug = true

	receiveChan := a.LaunchListener()

	for {
		select {
		case msg := <-receiveChan:
			switch m := msg.(type) {
			case *ast.AriConnected:
				infos, err := r.AsteriskInfoGet()
				if err != nil {
					fmt.Println("Couldn'get infos", err)
				} else {
					pretty.Printf("Remote Asterisk version: %s\n", infos.System.Version)
				}

				//lst, _ := r.SoundsGet("", "")
				//pretty.Printf("Sounds found: %# v\n", lst)

			case *ast.StasisStart:
				r.ChannelsPlayPostById(m.Channel.Id, rest.PlayParams{
					Media: "demo-moreinfo",
				})
			case *ast.StasisEnd:
				fmt.Println("Oh well, ended Stasis")
			case *ast.ChannelDtmfReceived:
				fmt.Println("Got DTMF:", m.Digit)
			case *ast.ChannelHangupRequest:
				fmt.Printf("Hangup for channel %s\n", m.Channel)
			case *ast.ChannelVarset:
				fmt.Printf("Setting channel variable: %s to '%s'\n", m.Variable, m.Value)
			default:
				pretty.Printf("Received some message: %# v\n", msg)
			}
		}

	}
}
