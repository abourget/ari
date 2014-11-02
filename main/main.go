package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/abourget/ari/models"
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
			case *models.AriConnected:
				infos, err := r.AsteriskInfoGet()
				if err != nil {
					fmt.Println("Couldn'get infos", err)
				} else {
					pretty.Printf("Remote Asterisk version: %s\n", infos.System.Version)
				}
			case *models.StasisStart:
				r.ChannelsPlayPostById(m.Channel.Id, rest.PlayParams{
					Media: "demo-congrats",
				})
			case *models.StasisEnd:
				fmt.Println("Oh well, ended Stasis")
			case *models.ChannelHangupRequest:
				fmt.Printf("Hangup for channel %s\n", m.Channel)
			case *models.ChannelVarset:
				fmt.Printf("Setting channel variable: %s to '%s'\n", m.Variable, m.Value)
			default:
				pretty.Printf("Received some message: %# v\n", msg)
			}
		}

	}
}
