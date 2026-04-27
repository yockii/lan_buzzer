package websocket

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type JoinPayload struct {
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
}

type BuzzPayload struct {
	Timestamp int64 `json:"timestamp"`
}

type StateChangedPayload struct {
	State     string `json:"state"`
	StartTime *int64 `json:"startTime,omitempty"`
}

type PlayerListPayload struct {
	Players []PlayerInfo `json:"players"`
}

type PlayerInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	DeviceType string `json:"deviceType"`
}

type BuzzResultPayload struct {
	Winner  *PlayerInfo `json:"winner,omitempty"`
	IsEarly bool        `json:"isEarly"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}
