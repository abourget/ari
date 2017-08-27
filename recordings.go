package ari

import "fmt"

//
// Service
//

type RecordingService struct {
	client *Client
}

func (s *RecordingService) ListStored() ([]*StoredRecording, error) {
	var out []*StoredRecording
	return out, s.client.Get("/recordings/stored", nil, &out)
}

func (s *RecordingService) GetStored(recordingName string) (*StoredRecording, error) {
	var out StoredRecording
	return &out, s.client.Get(fmt.Sprintf("/recordings/stored/%s", recordingName), nil, &out)
}

func (s *RecordingService) GetLive(recordingName string) (*LiveRecording, error) {
	var out LiveRecording
	return &out, s.client.Get(fmt.Sprintf("/recordings/live/%s", recordingName), nil, &out)
}

func (s *RecordingService) DeleteStored(recordingName string) error {
	return s.client.Delete(fmt.Sprintf("/recordings/stored/%s", recordingName), nil)
}

func (s *RecordingService) CopyStored(recordingName, destinationRecordingName string) (*StoredRecording, error) {
	var out StoredRecording
	params := map[string]string{
		"destinationRecordingName": destinationRecordingName,
	}

	return &out, s.client.Post(fmt.Sprintf("/recordings/stored/%s/copy", recordingName), params, &out)
}

//
// Models
//

type StoredRecording struct {
	Format string
	Name   string

	// For further manipulations
	client *Client
}

func (s *StoredRecording) setClient(client *Client) {
	s.client = client
}

func (s *StoredRecording) Delete() error {
	return s.client.Recordings.DeleteStored(s.Name)
}

func (s *StoredRecording) Copy(destinationRecordingName string) (*StoredRecording, error) {
	return s.client.Recordings.CopyStored(s.Name, destinationRecordingName)
}

type LiveRecording struct {
	Cause           string
	Duration        *int64
	Format          string
	Name            string
	SilenceDuration *int64 `json:"silence_duration"`
	State           string
	TalkingDuration *int64 `json:"talking_duration"`
	TargetURI       string `json:"target_uri"`

	// For further manipulations
	client *Client
}

func (l *LiveRecording) setClient(client *Client) {
	l.client = client
}

func (l *LiveRecording) Cancel() error {
	return l.client.Delete(fmt.Sprintf("/recordings/live/%s", l.Name), nil)
}

func (l *LiveRecording) Stop() error {
	return l.client.Post(fmt.Sprintf("/recordings/live/%s/stop", l.Name), nil, nil)
}

func (l *LiveRecording) Pause() error {
	return l.client.Post(fmt.Sprintf("/recordings/live/%s/pause", l.Name), nil, nil)
}

func (l *LiveRecording) Unpause() error {
	return l.client.Delete(fmt.Sprintf("/recordings/live/%s/pause", l.Name), nil)
}

func (l *LiveRecording) Mute() error {
	return l.client.Post(fmt.Sprintf("/recordings/live/%s/mute", l.Name), nil, nil)
}

func (l *LiveRecording) Unmute() error {
	return l.client.Delete(fmt.Sprintf("/recordings/live/%s/mute", l.Name), nil)
}
