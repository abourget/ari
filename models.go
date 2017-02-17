package ari

import "fmt"

type Variable struct {
	Value string
}

type CallerID struct {
	Name   string
	Number string
}

func (c *CallerID) String() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Number)
}

type Sound struct {
	Formats []FormatLangPair
	Id      string
	Text    string
}

type DialplanCEP struct {
	Context  string
	Exten    string
	Priority int
}

type Dialplan struct {
	Context  string `json:"context"`
	Exten    string `json:"extension"`
	Priority int    `json:"priority"`
	Label    string `json:"label"`
}

type FormatLangPair struct {
	Format   string
	Language string
}

type AsteriskInfo struct {
	Build  *BuildInfo
	Config *ConfigInfo
	Status *StatusInfo
	System *SystemInfo
}

type BuildInfo struct {
	Date    string
	Kernel  string
	Machine string
	Options string
	Os      string
	User    string
}

type ConfigInfo struct {
	DefaultLanguage string  `json:"default_language"`
	MaxChannels     int64   `json:"max_channels"`
	MaxLoad         float64 `json:"max_load"`
	MaxOpenFiles    int64   `json:"max_open_files"`
	Name            string
	SetId           SetId
}

type SetId struct {
	Group string
	User  string
}

type StatusInfo struct {
	LastReloadTime *Time `json:"last_reload_time"`
	StartupTime    *Time `json:"startup_time"`
}

type SystemInfo struct {
	EntityId string `json:"entity_id"`
	Version  string
}

type DeviceState struct {
	State string
	Name  string
}

type Endpoint struct {
	Technology string `json:"technology"`
	Resource   string `json:"resource"`
	State      string `json:"state"`
}

type Peer struct {
	PeerStatus string `json:"peer_status"`
}
