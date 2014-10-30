package ari

// Package ari implements the Asterisk ARI interface

import (
	"encoding/json"
	"fmt"

	"code.google.com/p/go.net/websocket"
)

type ARIClient struct {
	ws       *websocket.Conn
	hostname string
	login    string
	password string
	port     int
	ReceiveChan  chan interface{}
}

func NewARI(login, password, hostname string, port int) *ARIClient {
	ari := ARIClient{hostname: hostname, port: port, login: login, password: password}
	ari.ReceiveChan = make(chan interface{}, 100)
	return &ari
}
func (ari *ARIClient) Connect(appName string) error {
	ws, err := websocket.Dial(fmt.Sprintf("ws://%s:%d/ari/events?api_key=%s:%s&app=%s", ari.hostname, ari.port, ari.login, ari.password, appName), "", "http://localhost")
	ari.ws = ws
	return err
}

func (ari *ARIClient) HandleReceive() {
	for {
		fmt.Println("Listening using websocket.JSON.Receive...")
		var msg string
		websocket.Message.Receive(ari.ws, &msg)

		var data Message
		rawMsg := []byte(msg)
		err := json.Unmarshal(rawMsg, &data)
		if err != nil {
			fmt.Printf("Error decoding incoming '%#v': %s", msg, err)
			continue
		}

		//fmt.Printf("  -> %s", msg)

		msgType := data.Type
		var recvMsg interface{}
		switch msgType {
		case "StasisStart":
			recvMsg = &StasisStart{}
		case "StasisEnd":
			recvMsg = &StasisEnd{}
		case "ChannelVarset":
			recvMsg = &ChannelVarset{}
		default:
			recvMsg = &data
		}
		err = json.Unmarshal(rawMsg, recvMsg)
		if err != nil {
			fmt.Println("Error decoding structured message: %#v", err)
			continue
		}

		ari.ReceiveChan <- recvMsg
	}
}
