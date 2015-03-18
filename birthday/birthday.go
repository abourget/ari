package main

import (
	"fmt"
	"time"

	"github.com/abourget/ari"
)

type Birthday struct {
	client *ari.Client

	mixingBridge *ari.Bridge

	// Regarding outgoing channel
	outgoingChannel *ari.Channel
	currentPlayback *ari.Playback
	outgoingNumber  string

	// Regarding incoming channel
	incomingChannel *ari.Channel
	dtmfControlMode bool // Opposed to "Dial" mode
	dtmfDigits      string

	// Snooping stuff
	snoopingChannel  *ari.Channel
	currentRecording *ari.LiveRecording
}

func (b *Birthday) Setup() {
	// Create holding and mixing bridges
	mixingBridge, _ := b.client.Bridges.Create(ari.CreateBridgeParams{
		Type:     "mixing,proxy_media",
		BridgeId: "birthday-mixing",
		Name:     "Birthday App mixing bridge",
	})

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

func (b *Birthday) playSong() {
	b.stopOutgoingProcessing()

	playback, err := b.mixingBridge.Play(ari.PlayParams{Media: "sound:tt-weasels"}) // "sound:birthday-song"})

	if err != nil {
		b.incomingChannel.Play(ari.PlayParams{Media: "sound:invalid"})
		return
	}

	b.currentPlayback = playback
}

func (b *Birthday) playAndRecordSong() {
	b.stopOutgoingProcessing()

	b.playSong()
	b.recordSomething("song", false)
}

func (b *Birthday) recordMessage() {
	b.stopOutgoingProcessing()

	b.recordSomething("message", true)
}

func (b *Birthday) recordSomething(recordingType string, beep bool) {
	fmt.Println("***** Recording something")
	snoopChan, err := b.outgoingChannel.Snoop(ari.SnoopParams{
		Spy: "in",
		App: "birthday-snoop",
	})

	fmt.Println("***** Recording something: before snoopChan creation")
	if err != nil {
		fmt.Println("Snooping channel, couldn't create: ", err)
		return
	}

	fmt.Println("***** Recording something: snoopChan creation okay")
	b.snoopingChannel = snoopChan

	recording, err := b.snoopingChannel.Record(ari.RecordParams{
		Name:   fmt.Sprintf("birthday-%s-%s-%s", formattedNow(), b.outgoingNumber, recordingType),
		Format: "wav",
		Beep:   beep,
	})

	fmt.Println("***** Recording something: Record on snoopingChannel", err)

	if err != nil {
		b.incomingChannel.Play(ari.PlayParams{Media: "sound:invalid"})
		return
	}

	fmt.Println("***** Recording something: we're good, not playing error")

	b.currentRecording = recording
}

// StopOutgoingProcessing stops the playback and recordings on the outgoing channel.
func (b *Birthday) stopOutgoingProcessing() {
	if b.currentPlayback != nil {
		b.currentPlayback.Stop()
		b.currentPlayback = nil
	}
	if b.currentRecording != nil {
		b.currentRecording.Stop()
		b.currentRecording = nil
	}
}

func (b *Birthday) hangupOutgoing() {
	if b.outgoingChannel != nil {
		b.outgoingChannel.Hangup()
		b.outgoingChannel = nil
	}
}

func formattedNow() string {
	now := time.Now()
	return now.Format("2006-01-02-15-04-05")
}

func (b *Birthday) dialOutgoing(number string) {
	outgoing, err := b.client.Channels.Create(ari.OriginateParams{
		Endpoint: fmt.Sprintf("SIP/voipms/%s", b.outgoingNumber),
		App:      "birthday-outgoing",
	})
	if err != nil {
		if b.incomingChannel != nil {
			b.incomingChannel.Play(ari.PlayParams{Media: "sound:invalid"})
		}
		return
	}
	b.dtmfControlMode = true
	b.outgoingChannel = outgoing
}
