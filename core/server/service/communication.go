package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CommunicationEventMsg struct {
	Event   string `json:"event"`
	Message any    `json:"message"`
}

type CommunicationService interface {
	BroadcastToRoom(roomId string, msg CommunicationEventMsg) (*http.Response, error)
}

type communicationService struct {
	key    string
	client *http.Client
}

const (
	broadcastToRoomEndpoint = "/api/v1/room/%v/broadcast"
)

func NewCommunicationService(key string) CommunicationService {
	return &communicationService{
		key: key,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (cs communicationService) BroadcastToRoom(roomId string, msg CommunicationEventMsg) (*http.Response, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(msg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(broadcastToRoomEndpoint, roomId),
		&body,
	)
	if err != nil {
		return nil, err
	}

	res, err := cs.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
