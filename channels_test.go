package ari

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelString(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		expected string
		channel  Channel
	}{
		{"id=12345.12,state=Down",
			Channel{ID: "12345.12", State: "Down"}},
		{"id=12345.12,caller=Bob <911>,state=Up",
			Channel{ID: "12345.12", State: "Up", Caller: &CallerID{Name: "Bob", Number: "911"}}},
		{"id=12345.12,caller=Bob <911>,with=Alice <166>,state=Down",
			Channel{ID: "12345.12", State: "Down", Caller: &CallerID{Name: "Bob", Number: "911"}, Connected: &CallerID{Name: "Alice", Number: "166"}},
		},
	}

	for _, test := range tests {
		assert.Equal(test.expected, test.channel.String())
	}
}
