package main

import (
	"fmt"

	"github.com/abourget/ari"
	"github.com/kr/pretty"
)

type CallManager struct {
	client          *ari.Client
	liveRecording   *ari.LiveRecording
	currentPlayback *ari.Playback
	bridge          *ari.Bridge
	calledChannel   *ari.Channel
}

func (c *CallManager) handleMessage(msg interface{}) {
	switch m := msg.(type) {
	case *ari.AriConnected:
		fmt.Println("CallManager: ARI connected")

	case *ari.StasisStart:
		c.currentPlayback, _ = m.Channel.Play(ari.PlayParams{
			Media: "sound:demo-moreinfo",
		})
		m.Channel.SetVar("TALK_DETECT(set)", "")

	case *ari.StasisEnd:
		fmt.Println("Oh well, ended Stasis")

	case *ari.ChannelDtmfReceived:
		fmt.Println("Got DTMF:", m.Digit)

		if c.currentPlayback != nil {
			c.currentPlayback.Stop()
			c.currentPlayback = nil
		}

		if m.Digit == "1" {
			c.bridge, _ = c.client.Bridges.Create(ari.CreateBridgeParams{
				Type:     "mixing,proxy_media,dtmf_events",
				BridgeId: "mycall",
				Name:     "my-named-bridge",
			})

			otherChannel, err := c.client.Channels.Create(ari.OriginateParams{
				Endpoint: "SIP/voipms/5148888888",
				App:      "outgoing-call",
			})
			if err != nil {
				fmt.Println("Callmanager: error creating bridge:", err)
				return
			}

			c.calledChannel = otherChannel

			c.bridge.AddChannel(otherChannel.Id, ari.Participant)
		}

		if m.Digit == "2" {
			if c.calledChannel != nil {
				c.calledChannel.Hangup()
			}
		}

		if m.Digit == "*" {
			if c.liveRecording == nil {
				c.liveRecording, _ = m.Channel.Record(ari.RecordParams{
					Name:   "superbob",
					Format: "wav",
				})
			} else {
				c.liveRecording.Stop()
				c.liveRecording = nil
			}
		}

	case *ari.ChannelTalkingStarted:
		fmt.Println("CallManager: They started talking!")

	case *ari.ChannelTalkingFinished:
		fmt.Println("CallManager: They stopped talking! Talked for", m.Duration, "ms")

	case *ari.ChannelHangupRequest:
		fmt.Printf("CallManager: Hangup for channel %s\n", m.Channel)

	case *ari.ChannelVarset:
		fmt.Printf("CallManager: Setting channel variable: %s to '%s'\n", m.Variable, m.Value)

	case *ari.PlaybackStarted:
		fmt.Println("CallManager: Playback started")

	case *ari.PlaybackFinished:
		fmt.Println("CallManager: Playback finished: ", m.Playback.MediaURI)

	default:
		pretty.Printf("CallManager: Received some message: %+v\n", msg)
	}
}

func (c *CallManager) Listen() {
	receiveChan := c.client.LaunchListener()

	for {
		select {
		case msg := <-receiveChan:
			c.handleMessage(msg)
		}
	}
}
