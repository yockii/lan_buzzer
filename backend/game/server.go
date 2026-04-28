package game

import (
	"github.com/yockii/lan_qr/backend/quiz"
	"time"
)

func NewGameServer() *GameServer {
	return &GameServer{
		State:   StateWaiting,
		Players: make(map[string]*Player),
	}
}

func (gs *GameServer) AddPlayer(player *Player) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	gs.Players[player.ID] = player
}

func (gs *GameServer) RemovePlayer(playerID string) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	delete(gs.Players, playerID)
}

func (gs *GameServer) GetPlayers() []*Player {
	gs.Mutex.RLock()
	defer gs.Mutex.RUnlock()

	players := make([]*Player, 0, len(gs.Players))
	for _, player := range gs.Players {
		players = append(players, player)
	}
	return players
}

func (gs *GameServer) StartGame() {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	gs.State = StateReady
	now := time.Now()
	gs.StartTime = &now
	gs.WinnerID = nil
}

func (gs *GameServer) RecordBuzz(playerID string) bool {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	if gs.State != StateReady {
		return false
	}

	gs.State = StateLocked
	gs.WinnerID = &playerID
	return true
}

func (gs *GameServer) ResetGame() {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	gs.State = StateWaiting
	gs.WinnerID = nil
	gs.StartTime = nil
}

func (gs *GameServer) GetState() GameState {
	gs.Mutex.RLock()
	defer gs.Mutex.RUnlock()

	return gs.State
}

func (gs *GameServer) GetWinner() *Player {
	gs.Mutex.RLock()
	defer gs.Mutex.RUnlock()

	if gs.WinnerID == nil {
		return nil
	}

	return gs.Players[*gs.WinnerID]
}

// SetQuestionBank sets the question bank and enables quiz mode
func (gs *GameServer) SetQuestionBank(qb *quiz.QuestionBank) {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	gs.QuestionBank = qb
	gs.QuizState = &QuizState{
		Answers: make(map[string]*PlayerAnswer),
	}
	gs.Mode = ModeQuiz
}

// HasQuestionBank returns true if a question bank is set
func (gs *GameServer) HasQuestionBank() bool {
	gs.Mutex.RLock()
	defer gs.Mutex.RUnlock()
	return gs.QuestionBank != nil
}

// GetMode returns the current game mode
func (gs *GameServer) GetMode() GameMode {
	gs.Mutex.RLock()
	defer gs.Mutex.RUnlock()
	return gs.Mode
}

// StartQuizQuestion starts a new quiz question
func (gs *GameServer) StartQuizQuestion() *quiz.Question {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	if gs.QuestionBank == nil {
		return nil
	}

	// Mark current question as asked before starting a new one
	if gs.QuizState != nil && gs.QuizState.CurrentQuestion != nil {
		gs.QuestionBank.MarkAsked(gs.QuizState.CurrentQuestion.ID)
	}

	question := gs.QuestionBank.GetRandomQuestion()
	if question != nil {
		gs.QuizState.CurrentQuestion = question
		gs.QuizState.Answers = make(map[string]*PlayerAnswer)
		gs.QuizState.WinnerID = ""
	}

	return question
}

// SubmitAnswer submits a player's answer
func (gs *GameServer) SubmitAnswer(playerID string, answer string) *PlayerAnswer {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	if gs.QuizState == nil || gs.QuizState.CurrentQuestion == nil {
		return nil
	}

	question := gs.QuizState.CurrentQuestion

	// Create answer record
	playerAnswer := &PlayerAnswer{
		PlayerID:  playerID,
		Answer:    answer,
		Timestamp: time.Now().UnixMilli(),
	}

	// Auto-judge for single choice and true/false
	if !question.NeedsManualJudgment() {
		if question.IsCorrect(answer) {
			playerAnswer.Status = AnswerStatusCorrect
		} else {
			playerAnswer.Status = AnswerStatusIncorrect
		}
	} else {
		// Manual judgment for open-ended
		if question.IsCorrect(answer) {
			playerAnswer.Status = AnswerStatusCorrect
		} else {
			playerAnswer.Status = AnswerStatusPending
		}
	}

	gs.QuizState.Answers[playerID] = playerAnswer
	return playerAnswer
}

// JudgeAnswer manually judges a player's answer
func (gs *GameServer) JudgeAnswer(playerID string, correct bool) bool {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	if gs.QuizState == nil {
		return false
	}

	playerAnswer, exists := gs.QuizState.Answers[playerID]
	if !exists {
		return false
	}

	if correct {
		playerAnswer.Status = AnswerStatusCorrect
	} else {
		playerAnswer.Status = AnswerStatusIncorrect
		// If changing from correct to incorrect, clear IsWinner
		if playerAnswer.IsWinner {
			playerAnswer.IsWinner = false
		}
	}

	return true
}

// GetQuizState returns the current quiz state
func (gs *GameServer) GetQuizState() *QuizState {
	gs.Mutex.RLock()
	defer gs.Mutex.RUnlock()

	if gs.QuizState == nil {
		return nil
	}

	// Return a copy to avoid race conditions
	return &QuizState{
		CurrentQuestion: gs.QuizState.CurrentQuestion,
		Answers:        gs.QuizState.Answers,
		WinnerID:       gs.QuizState.WinnerID,
	}
}

// NextQuizQuestion moves to the next question
func (gs *GameServer) NextQuizQuestion() *quiz.Question {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	if gs.QuestionBank == nil {
		return nil
	}

	// Mark current question as asked
	if gs.QuizState != nil && gs.QuizState.CurrentQuestion != nil {
		gs.QuestionBank.MarkAsked(gs.QuizState.CurrentQuestion.ID)
	}

	// Get next question
	question := gs.QuestionBank.GetRandomQuestion()
	if question != nil {
		gs.QuizState.CurrentQuestion = question
		gs.QuizState.Answers = make(map[string]*PlayerAnswer)
	}

	return question
}

// ResetQuizQuestionBank resets the question bank to start over
func (gs *GameServer) ResetQuizQuestionBank() bool {
	gs.Mutex.Lock()
	defer gs.Mutex.Unlock()

	if gs.QuestionBank == nil {
		return false
	}

	gs.QuestionBank.ResetAsked()

	// Reset current quiz state
	if gs.QuizState != nil {
		gs.QuizState.CurrentQuestion = nil
		gs.QuizState.Answers = make(map[string]*PlayerAnswer)
		gs.QuizState.WinnerID = ""
	}

	return true
}
