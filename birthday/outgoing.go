package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

func (b *Birthday) handleOutgoingMessage(msg interface{}) {
	switch m := msg.(type) {
	case *ari.AriConnected:
		fmt.Println("Outgoing: ARI connected")

	case *ari.StasisStart:

		fmt.Println("Outgoing: Statis started, detecting speech")
		m.Channel.SetVar("TALK_DETECT(set)", "500")

		// Bridge with the other folk
		b.mixingBridge.AddChannel(m.Channel.ID, ari.Participant)

	case *ari.StasisEnd:
		fmt.Println("Outgoing: Statis ended")
	case *ari.ChannelDtmfReceived:
		fmt.Println("Outgoing: Got DTMF:", m.Digit)
	case *ari.ChannelTalkingStarted:
		fmt.Println("Outgoing: They started talking!")

	case *ari.ChannelTalkingFinished:
		fmt.Println("Outgoing: They stopped talking! Talked for", m.Duration, "ms")

	case *ari.ChannelDestroyed:
		fmt.Println("Outgoing: Channel destroyed, back to dialing")
		b.stopOutgoingProcessing()
		b.outgoingChannel = nil
		b.incomingChannel.Play(ari.PlayParams{
			Media: "sound:confbridge-has-left",
		})
		b.dtmfControlMode = true

	case *ari.ChannelHangupRequest:
		fmt.Printf("Outgoing: Hangup for channel %s\n", m.Channel)

	case *ari.ChannelVarset:
		//fmt.Printf("Outgoing: Setting channel variable: %s to '%s'\n", m.Variable, m.Value)

	case *ari.PlaybackStarted:
		fmt.Println("Outgoing: Playback started")

	case *ari.ChannelStateChange:
		fmt.Printf("Outgoing: ChannelStateChange: %#v\n", m.Channel)

	case *ari.PlaybackFinished:
		fmt.Println("Outgoing: Playback finished: ", m.Playback.MediaURI)

	default:
		pretty.Printf("Outgoing: Received some message: %+v\n", msg)
	}
}
