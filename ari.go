package ari

// Package ari implements the Asterisk ARI interface. See: https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+ARI

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/jmcvetta/napping"

	"golang.org/x/net/websocket"
)

type Client struct {
	Debug         bool
	ws            *websocket.Conn
	hostname      string
	username      string
	password      string
	port          int
	appName       string
	SubscribeAll  bool
	reconnections int

	session  napping.Session
	endpoint string

	// Services
	Channels     ChannelService
	Bridges      BridgeService
	Applications ApplicationService
	Asterisk     AsteriskService
	DeviceStates DeviceStateService
	Endpoints    EndpointService
	Events       EventService
	Mailboxes    MailboxService
	Playbacks    PlaybackService
	Recordings   RecordingService
	Sounds       SoundService
}

func NewClient(username, password, hostname string, port int, appName string) *Client {
	userinfo := url.UserPassword(username, password)
	endpoint := "http://" + net.JoinHostPort(hostname, strconv.Itoa(port))

	c := &Client{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		appName:  appName,
		session: napping.Session{
			Userinfo: userinfo,
		},
		endpoint: endpoint,
	}
	c.Channels.client = c
	c.Bridges.client = c
	c.Sounds.client = c
	c.Playbacks.client = c
	c.Asterisk.client = c
	c.Mailboxes.client = c
	c.Recordings.client = c
	c.Events.client = c
	c.Applications.client = c
	c.DeviceStates.client = c
	c.Endpoints.client = c

	return c
}

func (c *Client) LaunchListener() <-chan Eventer {
	ch := make(chan Eventer, 100)
	go c.handleReceive(ch)
	return ch
}

func (c *Client) handleReceive(ch chan<- Eventer) {
	for {
		c.reconnect(ch)
		c.listenForMessages(ch)
	}
}

func (c *Client) reconnect(ch chan<- Eventer) {
	for {
		err := c.connect()

		if err == nil {
			// Connected successfully
			fmt.Println("Connected to websocket successfully, registered", c.appName)
			ch <- &AriConnected{
				Reconnections: c.reconnections,
				Event:         Event{Message: Message{Type: "AriConnected"}},
			}
			c.reconnections++
			return
		}

		fmt.Println("Error connecting, trying in 3 seconds:", err)
		time.Sleep(3 * time.Second)
		continue
	}
}

func (c *Client) connect() error {
	url := fmt.Sprintf("ws://%s:%d/ari/events?api_key=%s:%s&app=%s&subscribeAll=%t", c.hostname, c.port, c.username, c.password, c.appName, c.SubscribeAll)
	ws, err := websocket.Dial(url, "", "http://localhost")
	c.ws = ws
	return err
}

func (c *Client) listenForMessages(ch chan<- Eventer) {
	for {
		var msg string
		err := websocket.Message.Receive(c.ws, &msg)
		if err != nil {
			fmt.Println("Whoops, error reading from Socket, resetting connection")
			ch <- &AriDisconnected{Event: Event{Message: Message{Type: "AriDisconnected"}}}
			return
		}

		if err != nil {
			fmt.Printf("Error decoding structured message: %#v\n", err)
			continue
		}

		//fmt.Printf("  -> %s", msg)
		recvMsg, err := parseMsg([]byte(msg))

		if err != nil {
			fmt.Printf("Error decoding incoming '%#v': %s\n", msg, err)
			continue
		}

		c.setClientRecurse(recvMsg)

		ch <- recvMsg
	}
}

func (c *Client) Log(format string, v ...interface{}) {
	if c.Debug {
		log.Println(c.appName, fmt.Sprintf(format, v...))
	}
}

func (c *Client) setClientRecurse(recvMsg interface{}) {
	original := reflect.ValueOf(recvMsg)
	doAssignClient(c, original, 0)
}

func doAssignClient(c *Client, original reflect.Value, depth int) {
	// based off: https://gist.github.com/hvoecking/10772475
	pkgPath := original.Type().PkgPath()

	if pkgPath == "time" {
		return
	}

	//fmt.Println("Ok, got something as a value, has PkgPath:", depth, original.Type().PkgPath(), original)

	if original.CanInterface() {
		setter, ok := original.Interface().(clientSetter)
		if ok {
			if !original.IsNil() {
				setter.setClient(c)
			}
			return
		}
	}

	switch original.Kind() {
	case reflect.Ptr:
		originalVal := original.Elem()
		if !originalVal.IsValid() {
			return
		}
		doAssignClient(c, originalVal, depth+1)
	//case reflect.Interface:
	//	originalVal := original.Interface()
	//	doAssignClient(c, originalVal)
	case reflect.Struct:
		for i := 0; i < original.NumField(); i++ {
			doAssignClient(c, original.Field(i), depth+1)
		}

	case reflect.Slice:
		for i := 0; i < original.Len(); i++ {
			doAssignClient(c, original.Index(i), depth+1)
		}
		//case reflect.Map:
		// we don't have that case in our model
		//default:
	}
}

type clientSetter interface {
	setClient(*Client)
}
