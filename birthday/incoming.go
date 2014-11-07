package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

func (b *Birthday) handleIncomingMessage(msg interface{}) {
	switch m := msg.(type) {
	case *ari.AriConnected:
		fmt.Println("Incoming: ARI connected")

	case *ari.StasisStart:
		fmt.Println("Incoming: Statis started, detecthing speech")
		b.incomingChannel = m.Channel

		b.currentPlayback, _ = m.Channel.Play(ari.PlayParams{
			Media: "sound:demo-moreinfo",
		})
		m.Channel.SetVar("TALK_DETECT(set)", "")

	case *ari.StasisEnd:
		fmt.Println("Incoming: Statis ended")

	case *ari.ChannelDtmfReceived:
		fmt.Println("Incoming: Got DTMF:", m.Digit)

		if b.currentPlayback != nil {
			b.currentPlayback.Stop()
			b.currentPlayback = nil
		}

		if m.Digit == "1" {
			b.bridge, _ = b.client.Bridges.Create(ari.CreateBridgeParams{
				Type:     "mixing,proxy_media,dtmf_events",
				BridgeId: "mycall",
				Name:     "my-named-bridge",
			})

			// otherChannel, err := c.client.Channels.Create(ari.OriginateParams{
			// 	Endpoint: "SIP/voipms/5149221144",
			// 	App:      "outgoing-call",
			// })
			// if err != nil {
			// 	fmt.Println("Incoming: error creating bridge:", err)
			// 	return
			// }

			// c.outgoingChannel = otherChannel

			b.bridge.AddChannel(m.Channel.Id, ari.Participant)
		}

		if m.Digit == "2" {
			if b.outgoingChannel != nil {
				b.outgoingChannel.Hangup()
			}
		}

		if m.Digit == "*" {
			if b.liveRecording == nil {
				b.liveRecording, _ = m.Channel.Record(ari.RecordParams{
					Name:   "superbob",
					Format: "wav",
				})
			} else {
				b.liveRecording.Stop()
				b.liveRecording = nil
			}
		}

	case *ari.ChannelTalkingStarted:
		fmt.Println("Incoming: They started talking!")

	case *ari.ChannelTalkingFinished:
		fmt.Println("Incoming: They stopped talking! Talked for", m.Duration, "ms")

	case *ari.ChannelHangupRequest:
		fmt.Printf("Incoming: Hangup for channel %s\n", m.Channel)

	case *ari.ChannelVarset:
		fmt.Printf("Incoming: Setting channel variable: %s to '%s'\n", m.Variable, m.Value)

	case *ari.PlaybackStarted:
		fmt.Println("Incoming: Playback started")

	case *ari.PlaybackFinished:
		fmt.Println("Incoming: Playback finished: ", m.Playback.MediaURI)

	default:
		pretty.Printf("Incoming: Received some message: %+v\n", msg)
	}
}
