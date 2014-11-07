package main

// Package main is a demo of feature and a sample program that uses `ari`, `rest` and `models`.

import "github.com/abourget/ari"

func main() {
	c := ari.NewClient("asterisk", "asterisk", "localhost", 8088, "birthday-incoming,birthday-outgoing")
	c.Debug = true

	birthday := &Birthday{
		client: c,
	}

	birthday.Setup()
	birthday.Listen()
}
