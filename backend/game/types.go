package game

import (
	"sync"
	"time"
)

type GameState string

const (
	StateWaiting GameState = "waiting"
	StateReady   GameState = "ready"
	StateLocked  GameState = "locked"
)

type Player struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	DeviceType  string    `json:"deviceType"`
	ConnectedAt time.Time `json:"connectedAt"`
}

type GameServer struct {
	State     GameState          `json:"state"`
	Players   map[string]*Player `json:"players"`
	WinnerID  *string            `json:"winnerId"`
	StartTime *time.Time         `json:"startTime"`
	Mutex     sync.RWMutex       `json:"-"`
}
