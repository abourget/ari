package ari

type EventService struct {
	client *Client
}

type Message struct {
	Type string
}

type Event struct {
	Message
	Application string
	Timestamp   *Time
}

func (e Event) GetApplication() string {
	return e.Application
}

func (e Event) GetType() string {
	return e.Type
}

type Eventer interface {
	GetApplication() string
	GetType() string
}
