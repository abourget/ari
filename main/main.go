package main

// Package main is a demo of feature and a sample program that uses `ari`, `rest` and `models`.

import "github.com/abourget/ari"

func main() {
	c := ari.NewClient("asterisk", "asterisk", "localhost", 8088, "incoming-agent")
	c.Debug = true

	manager := &CallManager{client: c}
	go manager.Listen()

	c = ari.NewClient("asterisk", "asterisk", "localhost", 8088, "outgoing-call")
	c.Debug = true

	outgoing := &Outgoing{client: c}
	outgoing.Listen()
}
