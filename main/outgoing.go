package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

type Outgoing struct {
	client          *ari.Client
	liveRecording   *ari.LiveRecording
	currentPlayback *ari.Playback
	bridge          *ari.Bridge
}

func (c *Outgoing) handleMessage(msg interface{}) {
	switch m := msg.(type) {
	case *ari.AriConnected:
		fmt.Println("Outgoing: ARI connected")

	case *ari.StasisStart:
		fmt.Println("Outgoing: Statis started")
		m.Channel.SetVar("TALK_DETECT(set)", "")
	case *ari.StasisEnd:
		fmt.Println("Outgoing: Statis ended")
	case *ari.ChannelDtmfReceived:
		fmt.Println("Outgoing: Got DTMF:", m.Digit)
	case *ari.ChannelTalkingStarted:
		fmt.Println("Outgoing: They started talking!")

	case *ari.ChannelTalkingFinished:
		fmt.Println("Outgoing: They stopped talking! Talked for", m.Duration, "ms")

	case *ari.ChannelHangupRequest:
		fmt.Printf("Outgoing: Hangup for channel %s\n", m.Channel)

	case *ari.ChannelVarset:
		fmt.Printf("Outgoing: Setting channel variable: %s to '%s'\n", m.Variable, m.Value)

	case *ari.PlaybackStarted:
		fmt.Println("Outgoing: Playback started")

	case *ari.ChannelStateChange:


	case *ari.PlaybackFinished:
		fmt.Println("Outgoing: Playback finished: ", m.Playback.MediaURI)

	default:
		pretty.Printf("Outgoing: Received some message: %+v\n", msg)
	}
}

func (o *Outgoing) Listen() {
	receiveChan := o.client.LaunchListener()

	for {
		select {
		case msg := <-receiveChan:
			o.handleMessage(msg)
		}
	}
}
