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

	if _, err := s.client.Get("/recordings/stored", nil, &out); err != nil {
		return nil, err
	}

	s.client.setClientRecurse(out)
	return out, nil
}

func (s *RecordingService) GetStored(recordingName string) (*StoredRecording, error) {
	var out StoredRecording

	if _, err := s.client.Get(fmt.Sprintf("/recordings/stored/%s", recordingName), nil, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

func (s *RecordingService) GetLive(recordingName string) (*LiveRecording, error) {
	var out LiveRecording

	if _, err := s.client.Get(fmt.Sprintf("/recordings/live/%s", recordingName), nil, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

func (s *RecordingService) DeleteStored(recordingName string) error {
	_, err := s.client.Delete(fmt.Sprintf("/recordings/stored/%s", recordingName), nil)
	return err
}

func (s *RecordingService) CopyStored(recordingName, destinationRecordingName string) (*StoredRecording, error) {
	var out StoredRecording
	params := map[string]string{
		"destinationRecordingName": destinationRecordingName,
	}

	if _, err := s.client.Post(fmt.Sprintf("/recordings/stored/%s/copy", recordingName), params, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
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
	_, err := l.client.Delete(fmt.Sprintf("/recordings/live/%s", l.Name), nil)
	return err
}

func (l *LiveRecording) Stop() error {
	_, err := l.client.Post(fmt.Sprintf("/recordings/live/%s/stop", l.Name), nil, nil)
	return err
}

func (l *LiveRecording) Pause() error {
	_, err := l.client.Post(fmt.Sprintf("/recordings/live/%s/pause", l.Name), nil, nil)
	return err
}

func (l *LiveRecording) Unpause() error {
	_, err := l.client.Delete(fmt.Sprintf("/recordings/live/%s/pause", l.Name), nil)
	return err
}

func (l *LiveRecording) Mute() error {
	_, err := l.client.Post(fmt.Sprintf("/recordings/live/%s/mute", l.Name), nil, nil)
	return err
}

func (l *LiveRecording) Unmute() error {
	_, err := l.client.Delete(fmt.Sprintf("/recordings/live/%s/mute", l.Name), nil)
	return err
}
