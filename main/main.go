package main

// Package main is a demo of feature and a sample program that uses `ari`, `rest` and `models`.

import (
	"fmt"
	"time"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

func main() {
	c := ari.NewClient("asterisk", "asterisk", "localhost", 8088, "hello-world")
	c.Debug = true

	receiveChan := c.LaunchListener()

	for {
		select {
		case msg := <-receiveChan:
			switch m := msg.(type) {
			case *ari.AriConnected:
				infos, err := c.Asterisk.GetInfo()
				if err != nil {
					fmt.Println("Couldn'get infos", err)
				} else {
					pretty.Printf("Remote Asterisk version: %s\n", infos.System.Version)
				}

				//lst, _ := r.SoundsGet("", "")
				//pretty.Printf("Sounds found: %# v\n", lst)

			case *ari.StasisStart:
				m.Channel.Play(ari.PlayParams{
					Media: "sound:demo-moreinfo",
				})
				m.Channel.SetVar("TALK_DETECT(set)", "")
			case *ari.StasisEnd:
				fmt.Println("Oh well, ended Stasis")
			case *ari.ChannelDtmfReceived:
				fmt.Println("Got DTMF:", m.Digit)
			case *ari.ChannelTalkingStarted:
				fmt.Println("They started talking!")
			case *ari.ChannelTalkingFinished:
				fmt.Println("They stopped talking! Talked for", m.Duration, "ms")
			case *ari.ChannelHangupRequest:
				fmt.Printf("Hangup for channel %s\n", m.Channel)
			case *ari.ChannelVarset:
				fmt.Printf("Setting channel variable: %s to '%s'\n", m.Variable, m.Value)
			case *ari.PlaybackStarted:
				go func() {
					time.Sleep(2 * time.Second)
					m.Playback.Stop()
				}()
			case *ari.PlaybackFinished:
				fmt.Println("Playback finished: ", m.Playback.MediaURI)
			default:
				pretty.Printf("Received some message: %# v\n", msg)
			}
		}

	}
}
