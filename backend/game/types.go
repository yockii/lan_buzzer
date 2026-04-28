package game

import (
	"github.com/yockii/lan_qr/backend/quiz"
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

	// Quiz mode fields
	QuestionBank *quiz.QuestionBank
	QuizState    *QuizState
	Mode         GameMode
}

// GameMode represents the current game mode
type GameMode string

const (
	ModeBuzzer GameMode = "buzzer"
	ModeQuiz   GameMode = "quiz"
)

// QuizState represents the state of a quiz session
type QuizState struct {
	CurrentQuestion *quiz.Question
	Answers         map[string]*PlayerAnswer
	WinnerID        string
	mutex           sync.RWMutex
}

// PlayerAnswer represents a player's answer to a quiz question
type PlayerAnswer struct {
	PlayerID  string
	Answer    string
	Timestamp int64
	Status    AnswerStatus
	IsWinner  bool
}

// AnswerStatus represents the status of a player's answer
type AnswerStatus string

const (
	AnswerStatusPending   AnswerStatus = "pending"
	AnswerStatusCorrect   AnswerStatus = "correct"
	AnswerStatusIncorrect AnswerStatus = "incorrect"
)

// Question type for quiz mode (simplified version)
type Question struct {
	ID       string
	Type     QuestionType
	Question string
	Options  []string
	Answer   string
}

// NeedsManualJudgment returns true if this question type requires manual judging
func (q *Question) NeedsManualJudgment() bool {
	return q.Type == OpenEnded
}

// IsCorrect checks if the given answer is correct
func (q *Question) IsCorrect(answer string) bool {
	return q.Answer == answer
}

// QuestionType represents the type of question
type QuestionType string

const (
	SingleChoice QuestionType = "single_choice"
	TrueFalse    QuestionType = "true_false"
	OpenEnded   QuestionType = "open_ended"
)

// QuestionBank manages quiz questions (interface)
type QuestionBank interface {
	HasQuestions() bool
	GetRandomQuestion() *Question
	MarkAsked(id string)
	ResetAsked()
}

