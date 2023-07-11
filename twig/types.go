package twig

import "encoding/json"

type GenericTwigMessage struct {
	Op    int8            `json:"op"`
	Event string          `json:"e,omitempty"`
	Data  json.RawMessage `json:"d,omitempty"`
}

type HolaMessage struct {
	Op   int8 `json:"op"`
	Data struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

type AuthMessage struct {
	Op   int8 `json:"op"`
	Data struct {
		Token string `json:"token"`
	} `json:"d"`
}

type CreateLogMessage struct {
	Op   int8          `json:"op"`
	Data CreateLogData `json:"d"`
}

type CreateLogData struct {
	Message          string `json:"message"`
	Level            string `json:"level"`
	ProjectNamespace string `json:"project_namespace"`
	FeedName         string `json:"feed_name"`
}
