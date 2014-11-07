package main

import "github.com/abourget/ari"

type Birthday struct {
	client *ari.Client

	// Regarding outgoing channel
	liveRecording   *ari.LiveRecording
	currentPlayback *ari.Playback
	holdingBridge   *ari.Bridge
	mixingBridge    *ari.Bridge

	// Regarding incoming channel
	outgoingChannel *ari.Channel
	incomingChannel *ari.Channel
}

func (b *Birthday) Setup() {
	// Create holding and mixing bridges
	holdingBridge, _ := b.client.Bridges.Create(ari.CreateBridgeParams{
		Type:     "holding",
		BridgeId: "birthday-holding",
		Name:     "Birthday App holding bridge",
	})

	mixingBridge, _ := b.client.Bridges.Create(ari.CreateBridgeParams{
		Type:     "mixing,proxy_media,dtmf_events",
		BridgeId: "birthday-mixing",
		Name:     "Birthday App mixing bridge",
	})

	b.holdingBridge = holdingBridge
	b.mixingBridge = mixingBridge
}

func (b *Birthday) Listen() {
	receiveChan := b.client.LaunchListener()

	for {
		select {
		case msg := <-receiveChan:
			appName := msg.GetApplication()
			if appName == "birthday-outgoing" || appName == "" {
				b.handleOutgoingMessage(msg)
			}
			if appName == "birthday-incoming" || appName == "" {
				b.handleIncomingMessage(msg)
			}
		}
	}

}
