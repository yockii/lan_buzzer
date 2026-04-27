package game

import (
	"sync"
	"time"

	"github.com/yockii/lan_qr/backend/quiz"
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

	// Quiz mode fields
	QuestionBank *quiz.QuestionBank `json:"-"`
	QuizState    *QuizState         `json:"-"`
	Mode         GameMode           `json:"mode"`
}

// GameMode represents the current game mode
type GameMode string

const (
	ModeBuzzer GameMode = "buzzer"
	ModeQuiz   GameMode = "quiz"
)

// QuizState represents the state of a quiz session
type QuizState struct {
	CurrentQuestion *quiz.Question         `json:"currentQuestion"`
	Answers         map[string]*PlayerAnswer `json:"answers"`
	WinnerID        string                 `json:"winnerId"`
	mutex           sync.RWMutex          `json:"-"`
}

// PlayerAnswer represents a player's answer to a quiz question
type PlayerAnswer struct {
	PlayerID  string       `json:"playerId"`
	Answer    string       `json:"answer"`
	Timestamp int64        `json:"timestamp"`
	Status    AnswerStatus `json:"status"`
	IsWinner  bool         `json:"isWinner"`
}

// AnswerStatus represents the status of a player's answer
type AnswerStatus string

const (
	AnswerStatusPending   AnswerStatus = "pending"
	AnswerStatusCorrect   AnswerStatus = "correct"
	AnswerStatusIncorrect AnswerStatus = "incorrect"
)
