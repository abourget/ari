package ari

// Package ari implements the Asterisk ARI interface. See: https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+ARI

import (
	"encoding/json"
	"fmt"
	"time"

	ast "github.com/abourget/ari/models"
	"github.com/abourget/ari/rest"

	"code.google.com/p/go.net/websocket"
)

type ARIClient struct {
	Debug bool
	ws            *websocket.Conn
	hostname      string
	username      string
	password      string
	port          int
	appName       string
	reconnections int
}

func NewARI(username, password, hostname string, port int, appName string) *ARIClient {
	ari := ARIClient{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		appName:  appName,
	}
	return &ari
}

func (ari *ARIClient) GetREST() *rest.REST {
	endpoint := fmt.Sprintf("http://%s:%d", ari.hostname, ari.port)
	r := rest.New(endpoint, ari.username, ari.password)
	return r
}

func (ari *ARIClient) LaunchListener() <-chan interface{} {
	ch := make(chan interface{}, 100)
	go ari.handleReceive(ch)
	return ch
}
func (ari *ARIClient) handleReceive(ch chan<- interface{}) {
	for {
		ari.reconnect(ch)
		ari.listenForMessages(ch)
	}
}

func (ari *ARIClient) reconnect(ch chan<- interface{}) {
	for {
		err := ari.connect()

		if err == nil {
			// Connected successfully
			fmt.Println("Connected to websocket successfully")
			ch <- &ast.AriConnected{ari.reconnections}
			ari.reconnections += 1
			return
		}

		fmt.Println("Error connecting, trying in 3 seconds:", err)
		time.Sleep(3 * time.Second)
		continue
	}
}

func (ari *ARIClient) connect() error {
	url := fmt.Sprintf("ws://%s:%d/ari/events?api_key=%s:%s&app=%s", ari.hostname, ari.port, ari.username, ari.password, ari.appName)
	ws, err := websocket.Dial(url, "", "http://localhost")
	ari.ws = ws
	return err
}

func (ari *ARIClient) listenForMessages(ch chan<- interface{}) {
	for {
		var msg string
		err := websocket.Message.Receive(ari.ws, &msg)
		if err != nil {
			fmt.Println("Whoops, error reading from Socket, resetting connection")
			ch <- &ast.AriDisconnected{}
			return
		}

		var data ast.Message
		rawMsg := []byte(msg)
		err = json.Unmarshal(rawMsg, &data)
		if err != nil {
			fmt.Printf("Error decoding incoming '%#v': %s", msg, err)
			continue
		}

		//fmt.Printf("  -> %s", msg)

		msgType := data.Type
		var recvMsg interface{}
		switch msgType {
		case "ChannelVarset":
			recvMsg = &ast.ChannelVarset{}
		case "ChannelDtmfReceived":
			recvMsg = &ast.ChannelDtmfReceived{}
		case "ChannelHangupRequest":
			recvMsg = &ast.ChannelHangupRequest{}
		case "StasisStart":
			recvMsg = &ast.StasisStart{}
		case "PlaybackStarted":
			recvMsg = &ast.PlaybackStarted{}
		case "PlaybackFinished":
			recvMsg = &ast.PlaybackFinished{}
		case "ChannelTalkingStarted":
			recvMsg = &ast.ChannelTalkingStarted{}
		case "ChannelTalkingFinished":
			recvMsg = &ast.ChannelTalkingFinished{}
		case "ChannelDialplan":
			recvMsg = &ast.ChannelDialplan{}
		case "ChannelCreated":
			recvMsg = &ast.ChannelCreated{}
		case "ChannelDestroyed":
			recvMsg = &ast.ChannelDestroyed{}
		case "StasisEnd":
			recvMsg = &ast.StasisEnd{}
		default:
			recvMsg = &data
		}
		err = json.Unmarshal(rawMsg, recvMsg)
		if err != nil {
			fmt.Println("Error decoding structured message: %#v", err)
			continue
		}

		ch <- recvMsg
	}
}
