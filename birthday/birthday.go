package main

import "github.com/abourget/ari"

type Birthday struct {
	client *ari.Client

	// Regarding outgoing channel
	liveRecording   *ari.LiveRecording
	currentPlayback *ari.Playback
	bridge          *ari.Bridge

	// Regarding incoming channel
	outgoingChannel *ari.Channel
	incomingChannel *ari.Channel
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
