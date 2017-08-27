package ari

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetClientRecurse(t *testing.T) {
	assert := assert.New(t)
	client := &Client{}

	// test nested struct
	msg := ChannelEnteredBridge{
		Bridge:  &Bridge{},
		Channel: &Channel{},
	}

	assert.Nil(msg.Bridge.client)
	assert.Nil(msg.Channel.client)

	client.setClientRecurse(&msg)
	assert.Equal(client, msg.Bridge.client)
	assert.Equal(client, msg.Channel.client)

	// test slice
	slice := []*Channel{&Channel{}}
	client.setClientRecurse(&slice)
	assert.Equal(client, slice[0].client)
}
