package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

func (b *Birthday) handleIncomingMessage(msg interface{}) {
	switch m := msg.(type) {
	case *ari.AriConnected:
		//fmt.Println("Incoming: ARI connected, setting up mixing bridge")
		b.Setup()

	case *ari.StasisStart:
		//fmt.Println("Incoming: Statis started, detecting speech")
		b.incomingChannel = m.Channel

		b.currentPlayback, _ = m.Channel.Play(ari.PlayParams{
			Media: "sound:hello-world",
		})

		b.mixingBridge.AddChannel(m.Channel.Id, ari.Participant)

	case *ari.StasisEnd:
		fmt.Println("Incoming: Statis ended")

	case *ari.ChannelDtmfReceived:
		fmt.Println("Incoming: Got DTMF:", b.dtmfDigits, m.Digit)

		if b.dtmfControlMode {
			if m.Digit == "*" {
				b.hangupOutgoing()
				b.dtmfControlMode = false
			} else if m.Digit == "1" {
				b.playSong()
			} else if m.Digit == "4" {
				b.playAndRecordSong()
			} else if m.Digit == "7" {
				b.recordMessage()
			} else if m.Digit == "#" {
				b.stopOutgoingProcessing()
			}
		} else {
			// Dial mode
			if m.Digit != "#" && m.Digit != "*" {
				b.dtmfDigits = fmt.Sprintf("%s%s", b.dtmfDigits, m.Digit)
			} else if m.Digit == "*" {
				b.dtmfDigits = ""
			} else if m.Digit == "#" {
				// Do the dial, and reset digits
				b.outgoingNumber = b.dtmfDigits
				b.dialOutgoing(b.outgoingNumber)
				b.dtmfDigits = ""
			}
		}

	case *ari.ChannelTalkingStarted:
		//fmt.Println("Incoming: They started talking!")

	case *ari.ChannelTalkingFinished:
		//fmt.Println("Incoming: They stopped talking! Talked for", m.Duration, "ms")

	case *ari.ChannelHangupRequest:
		fmt.Printf("Incoming: Hangup for channel %s\n", m.Channel)
		b.hangupOutgoing()

	case *ari.ChannelVarset:
		//fmt.Printf("Incoming: Setting channel variable: %s to '%s'\n", m.Variable, m.Value)

	case *ari.PlaybackStarted:
		//fmt.Println("Incoming: Playback started")

	case *ari.PlaybackFinished:
		//fmt.Println("Incoming: Playback finished: ", m.Playback.MediaURI)

	default:
		pretty.Printf("Incoming: Received some message: %+v\n", msg)
	}
}
