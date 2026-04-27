# Question Bank Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add quiz mode to LAN Buzzer where players compete to answer questions from a question bank, supporting multiple choice, true/false, and open-ended questions with automatic and manual judging.

**Architecture:** Backend quiz package manages question loading and random selection. WebSocket handles real-time answer broadcasting. Frontend components display questions and collect answers. Mode auto-detects based on presence of questions.txt file.

**Tech Stack:** Go 1.25+, Vue 3, TypeScript, TailwindCSS, WebSocket, Fiber v2

---

## Phase 1: Backend - Question Data Structures

### Task 1: Create Question Data Structures

**Files:**
- Create: `backend/quiz/question.go`
- Test: `backend/quiz/question_test.go`

- [ ] **Step 1: Write the failing test**

```go
package quiz

import "testing"

func TestQuestionTypes(t *testing.T) {
    mcQ := &Question{
        ID:       "q1",
        Type:     "single_choice",
        Question: "What is CPU?",
        Options:  []string{"A. CPU", "B. GPU", "C. RAM", "D. Disk"},
        Answer:   "A",
    }

    if mcQ.Type != "single_choice" {
        t.Errorf("Expected single_choice, got %s", mcQ.Type)
    }

    tfQ := &Question{
        ID:       "q2",
        Type:     "true_false",
        Question: "Earth is round?",
        Options:  []string{"对", "错"},
        Answer:   "对",
    }

    if tfQ.Type != "true_false" {
        t.Errorf("Expected true_false, got %s", tfQ.Type)
    }

    oeQ := &Question{
        ID:       "q3",
        Type:     "open_ended",
        Question: "Capital of China?",
        Answer:   "北京",
    }

    if oeQ.Type != "open_ended" {
        t.Errorf("Expected open_ended, got %s", oeQ.Type)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend/quiz && go test -v`
Expected: FAIL with "undefined: Question"

- [ ] **Step 3: Write minimal implementation**

```go
package quiz

// QuestionType represents the type of question
type QuestionType string

const (
    SingleChoice QuestionType = "single_choice"
    TrueFalse    QuestionType = "true_false"
    OpenEnded   QuestionType = "open_ended"
)

// Question represents a quiz question
type Question struct {
    ID       string
    Type     QuestionType
    Question string
    Options  []string // For single_choice and true_false
    Answer   string   // Correct answer
}

// IsCorrect checks if the given answer is correct
func (q *Question) IsCorrect(answer string) bool {
    return q.Answer == answer
}

// NeedsManualJudgment returns true if this question type requires manual judging
func (q *Question) NeedsManualJudgment() bool {
    return q.Type == OpenEnded
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend/quiz && go test -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/quiz/question.go backend/quiz/question_test.go
git commit -m "feat(backend): add question data structures"
```

### Task 2: Create Question Parser

**Files:**
- Create: `backend/quiz/parser.go`
- Test: `backend/quiz/parser_test.go`
- Test fixture: `backend/quiz/test_questions.txt`

- [ ] **Step 1: Create test fixture file**

Create `backend/quiz/test_questions.txt`:
```
[单选]
什么是CPU?|A.中央处理器|B.显卡|C.内存|D.硬盘|A

[判断]
地球是圆的|对

[问答]
中国的首都在哪里？|北京
```

- [ ] **Step 2: Write the failing test**

```go
package quiz

import (
    "os"
    "testing"
)

func TestParseQuestions(t *testing.T) {
    file, err := os.Open("test_questions.txt")
    if err != nil {
        t.Fatalf("Failed to open test file: %v", err)
    }
    defer file.Close()

    questions, err := ParseQuestions(file)
    if err != nil {
        t.Fatalf("Failed to parse questions: %v", err)
    }

    if len(questions) != 3 {
        t.Errorf("Expected 3 questions, got %d", len(questions))
    }

    // Check first question (single choice)
    q1 := questions[0]
    if q1.Type != SingleChoice {
        t.Errorf("Expected single_choice, got %s", q1.Type)
    }
    if q1.Question != "什么是CPU?" {
        t.Errorf("Expected '什么是CPU?', got '%s'", q1.Question)
    }
    if len(q1.Options) != 4 {
        t.Errorf("Expected 4 options, got %d", len(q1.Options))
    }
    if q1.Answer != "A" {
        t.Errorf("Expected answer 'A', got '%s'", q1.Answer)
    }

    // Check second question (true/false)
    q2 := questions[1]
    if q2.Type != TrueFalse {
        t.Errorf("Expected true_false, got %s", q2.Type)
    }
    if len(q2.Options) != 2 {
        t.Errorf("Expected 2 options, got %d", len(q2.Options))
    }

    // Check third question (open-ended)
    q3 := questions[2]
    if q3.Type != OpenEnded {
        t.Errorf("Expected open_ended, got %s", q3.Type)
    }
    if len(q3.Options) != 0 {
        t.Errorf("Expected 0 options for open-ended, got %d", len(q3.Options))
    }
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd backend/quiz && go test -v`
Expected: FAIL with "undefined: ParseQuestions"

- [ ] **Step 4: Write minimal implementation**

```go
package quiz

import (
    "bufio"
    "fmt"
    "io"
    "strings"
)

// ParseQuestions parses questions from a text file
func ParseQuestions(file io.Reader) ([]*Question, error) {
    scanner := bufio.NewScanner(file)
    var questions []*Question
    var currentType QuestionType
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        line := strings.TrimSpace(scanner.Text())

        // Skip empty lines
        if line == "" {
            continue
        }

        // Check for type header
        if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
            typeStr := strings.Trim(line, "[]")
            switch typeStr {
            case "单选":
                currentType = SingleChoice
            case "判断":
                currentType = TrueFalse
            case "问答":
                currentType = OpenEnded
            default:
                return nil, fmt.Errorf("line %d: unknown question type '%s'", lineNum, typeStr)
            }
            continue
        }

        // Parse question
        parts := strings.Split(line, "|")
        if len(parts) < 2 {
            return nil, fmt.Errorf("line %d: invalid format, expected at least 2 fields separated by |", lineNum)
        }

        question := &Question{
            ID:       fmt.Sprintf("q%d", lineNum),
            Type:     currentType,
            Question: parts[0],
            Answer:   parts[len(parts)-1],
        }

        // Add options for single_choice and true_false
        if currentType == SingleChoice || currentType == TrueFalse {
            // Parts[1:-1] are options
            question.Options = parts[1 : len(parts)-1]
        }

        questions = append(questions, question)
    }

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading file: %w", err)
    }

    return questions, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd backend/quiz && go test -v`
Expected: PASS

- [ ] **Step 6: Test with invalid format**

Add to `parser_test.go`:
```go
func TestParseQuestionsInvalidFormat(t *testing.T) {
    input := "invalid question without separator"
    reader := strings.NewReader(input)

    _, err := ParseQuestions(reader)
    if err == nil {
        t.Error("Expected error for invalid format, got nil")
    }
}
```

Add import: `"strings"`

Run: `cd backend/quiz && go test -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add backend/quiz/parser.go backend/quiz/parser_test.go backend/quiz/test_questions.txt
git commit -m "feat(backend): add question parser for txt format"
```

### Task 3: Create Question Bank

**Files:**
- Create: `backend/quiz/question_bank.go`
- Test: `backend/quiz/question_bank_test.go`

- [ ] **Step 1: Write the failing test**

```go
package quiz

import (
    "math/rand"
    "testing"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func TestQuestionBank(t *testing.T) {
    questions := []*Question{
        {ID: "q1", Type: SingleChoice, Question: "Q1", Answer: "A"},
        {ID: "q2", Type: SingleChoice, Question: "Q2", Answer: "B"},
        {ID: "q3", Type: SingleChoice, Question: "Q3", Answer: "C"},
    }

    qb := NewQuestionBank(questions)

    // Test HasQuestions
    if !qb.HasQuestions() {
        t.Error("Expected HasQuestions to return true")
    }

    // Test GetRandomQuestion
    q := qb.GetRandomQuestion()
    if q == nil {
        t.Error("Expected question, got nil")
    }

    // Test MarkAsked
    qb.MarkAsked(q.ID)
    q2 := qb.GetRandomQuestion()
    if q2 == nil {
        t.Error("Expected question, got nil")
    }

    // Eventually all questions should be asked
    for i := 0; i < 10; i++ {
        q := qb.GetRandomQuestion()
        if q != nil {
            qb.MarkAsked(q.ID)
        }
    }

    // After resetting, should get questions again
    qb.ResetAsked()
    q3 := qb.GetRandomQuestion()
    if q3 == nil {
        t.Error("Expected question after reset, got nil")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend/quiz && go test -v`
Expected: FAIL with "undefined: NewQuestionBank"

- [ ] **Step 3: Write minimal implementation**

```go
package quiz

import (
    "math/rand"
    "time"
)

// QuestionBank manages a collection of questions
type QuestionBank struct {
    questions []*Question
    askedIDs  map[string]bool
}

// NewQuestionBank creates a new question bank
func NewQuestionBank(questions []*Question) *QuestionBank {
    return &QuestionBank{
        questions: questions,
        askedIDs:  make(map[string]bool),
    }
}

// HasQuestions returns true if there are questions available
func (qb *QuestionBank) HasQuestions() bool {
    return len(qb.questions) > 0
}

// GetRandomQuestion returns a random unasked question
func (qb *QuestionBank) GetRandomQuestion() *Question {
    // Filter unasked questions
    var unasked []*Question
    for _, q := range qb.questions {
        if !qb.askedIDs[q.ID] {
            unasked = append(unasked, q)
        }
    }

    // If all questions asked, reset
    if len(unasked) == 0 && len(qb.questions) > 0 {
        qb.ResetAsked()
        unasked = qb.questions
    }

    // Still no questions (empty bank)
    if len(unasked) == 0 {
        return nil
    }

    // Return random question
    rand.Seed(time.Now().UnixNano())
    idx := rand.Intn(len(unasked))
    return unasked[idx]
}

// MarkAsked marks a question as asked
func (qb *QuestionBank) MarkAsked(id string) {
    qb.askedIDs[id] = true
}

// ResetAsked clears the asked list
func (qb *QuestionBank) ResetAsked() {
    qb.askedIDs = make(map[string]bool)
}

// LoadFromFile loads questions from a file
func LoadQuestionBank(filePath string) (*QuestionBank, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    questions, err := ParseQuestions(file)
    if err != nil {
        return nil, err
    }

    return NewQuestionBank(questions), nil
}
```

Add import: `"os"`

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend/quiz && go test -v`
Expected: PASS

- [ ] **Step 5: Test deduplication**

Add to `question_bank_test.go`:
```go
func TestQuestionBankDeduplication(t *testing.T) {
    questions := []*Question{
        {ID: "q1", Type: SingleChoice, Question: "Q1", Answer: "A"},
        {ID: "q2", Type: SingleChoice, Question: "Q2", Answer: "B"},
    }

    qb := NewQuestionBank(questions)

    // Get same question twice before marking
    q1 := qb.GetRandomQuestion()
    q2 := qb.GetRandomQuestion()

    // Without marking, should get questions
    if q1 == nil || q2 == nil {
        t.Error("Expected questions")
    }

    // Mark first as asked
    qb.MarkAsked(q1.ID)

    // Should only get the second one now
    q3 := qb.GetRandomQuestion()
    if q3.ID != q2.ID {
        t.Errorf("Expected only unasked question, got %s", q3.ID)
    }
}
```

Run: `cd backend/quiz && go test -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add backend/quiz/question_bank.go backend/quiz/question_bank_test.go
git commit -m "feat(backend): add question bank with deduplication"
```

### Task 4: Extend Game Server for Quiz Mode

**Files:**
- Modify: `backend/game/server.go`
- Modify: `backend/game/types.go`

- [ ] **Step 1: Add quiz state structures to types.go**

Read `backend/game/types.go` first to see existing structures, then add:

```go
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
```

Add imports if needed:
```go
import (
    "github.com/yockii/lan_qr/backend/quiz"
    "sync"
)
```

- [ ] **Step 2: Add quiz state to GameServer in server.go**

Read `backend/game/server.go` to see the GameServer struct, then add field:

```go
type GameServer struct {
    state     GameState
    players   []*Player
    buzzes    []Buzz
    mutex     sync.RWMutex
    StartTime *time.Time

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
```

- [ ] **Step 3: Add quiz methods to GameServer**

Add to `backend/game/server.go`:

```go
// SetQuestionBank sets the question bank and enables quiz mode
func (gs *GameServer) SetQuestionBank(qb *quiz.QuestionBank) {
    gs.mutex.Lock()
    defer gs.mutex.Unlock()

    gs.QuestionBank = qb
    gs.QuizState = &QuizState{
        Answers: make(map[string]*PlayerAnswer),
    }
    gs.Mode = ModeQuiz
}

// GetMode returns the current game mode
func (gs *GameServer) GetMode() GameMode {
    gs.mutex.RLock()
    defer gs.mutex.RUnlock()
    return gs.Mode
}

// StartQuizQuestion starts a new quiz question
func (gs *GameServer) StartQuizQuestion() *quiz.Question {
    gs.mutex.Lock()
    defer gs.mutex.Unlock()

    if gs.QuestionBank == nil {
        return nil
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
    gs.mutex.Lock()
    defer gs.mutex.Unlock()

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
            // Check if first correct answer
            if gs.QuizState.WinnerID == "" {
                playerAnswer.IsWinner = true
                gs.QuizState.WinnerID = playerID
            }
        } else {
            playerAnswer.Status = AnswerStatusIncorrect
        }
    } else {
        // Manual judgment for open-ended
        if question.IsCorrect(answer) {
            playerAnswer.Status = AnswerStatusCorrect
            if gs.QuizState.WinnerID == "" {
                playerAnswer.IsWinner = true
                gs.QuizState.WinnerID = playerID
            }
        } else {
            playerAnswer.Status = AnswerStatusPending
        }
    }

    gs.QuizState.Answers[playerID] = playerAnswer
    return playerAnswer
}

// JudgeAnswer manually judges a player's answer
func (gs *GameServer) JudgeAnswer(playerID string, correct bool) bool {
    gs.mutex.Lock()
    defer gs.mutex.Unlock()

    if gs.QuizState == nil {
        return false
    }

    playerAnswer, exists := gs.QuizState.Answers[playerID]
    if !exists {
        return false
    }

    if correct {
        playerAnswer.Status = AnswerStatusCorrect
        // Check if first correct answer
        if gs.QuizState.WinnerID == "" {
            playerAnswer.IsWinner = true
            gs.QuizState.WinnerID = playerID
            return true
        }
    } else {
        playerAnswer.Status = AnswerStatusIncorrect
    }

    return false
}

// GetQuizState returns the current quiz state
func (gs *GameServer) GetQuizState() *QuizState {
    gs.mutex.RLock()
    defer gs.mutex.RUnlock()

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
    gs.mutex.Lock()
    defer gs.mutex.Unlock()

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
        gs.QuizState.WinnerID = ""
    }

    return question
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/game/server.go backend/game/types.go
git commit -m "feat(backend): extend game server for quiz mode"
```

### Task 5: Integrate Question Bank Loading

**Files:**
- Modify: `backend/main.go`

- [ ] **Step 1: Load question bank on startup**

Read `backend/main.go` to find where GameServer is created, then add:

```go
gameServer := game.NewGameServer()

// Try to load question bank
questionBankPath := "questions.txt"
if qb, err := quiz.LoadQuestionBank(questionBankPath); err == nil {
    gameServer.SetQuestionBank(qb)
    log.Printf("Loaded %d questions from %s, mode: quiz", len(qb.HasQuestions()), questionBankPath)
} else {
    log.Printf("No question bank found (%v), mode: buzzer", err)
}
```

Add import:
```go
import "github.com/yockii/lan_qr/backend/quiz"
```

Note: Fix the log statement - HasQuestions is a method:
```go
log.Printf("Loaded question bank from %s, mode: quiz", questionBankPath)
```

- [ ] **Step 2: Build and test**

Run: `cd backend && go build`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add backend/main.go
git commit -m "feat(backend): load question bank on startup"
```

### Task 6: Add Quiz WebSocket Messages

**Files:**
- Modify: `backend/websocket/handler.go`

- [ ] **Step 1: Add quiz message handlers**

Read `backend/websocket/handler.go` to see existing handlers, then add to the switch statement:

```go
case "quiz_start":
    h.handleQuizStart(conn)
case "quiz_answer":
    h.handleQuizAnswer(conn, playerID, msg)
case "quiz_judge":
    h.handleQuizJudge(conn, msg)
case "quiz_next":
    h.handleQuizNext(conn)
```

- [ ] **Step 2: Implement quiz message handlers**

Add to `backend/websocket/handler.go`:

```go
func (h *Handler) handleQuizStart(conn *websocket.Conn) {
    question := h.server.StartQuizQuestion()
    if question == nil {
        h.sendError(conn, "No questions available")
        return
    }

    h.broadcastQuizQuestion(question)
}

func (h *Handler) handleQuizAnswer(conn *websocket.Conn, playerID string, msg *Message) {
    payload, ok := msg.Payload.(map[string]any)
    if !ok {
        h.sendError(conn, "Invalid payload")
        return
    }

    answer, _ := payload["answer"].(string)
    if answer == "" {
        return
    }

    playerAnswer := h.server.SubmitAnswer(playerID, answer)
    if playerAnswer != nil {
        h.broadcastQuizAnswerUpdate(playerAnswer)
    }
}

func (h *Handler) handleQuizJudge(conn *websocket.Conn, msg *Message) {
    payload, ok := msg.Payload.(map[string]any)
    if !ok {
        return
    }

    playerID, _ := payload["playerId"].(string)
    correctStr, _ := payload["correct"].(string)
    if playerID == "" || correctStr == "" {
        return
    }

    correct := correctStr == "true" || correctStr == "1"
    isFirstCorrect := h.server.JudgeAnswer(playerID, correct)

    if isFirstCorrect {
        // Broadcast the updated answer
        quizState := h.server.GetQuizState()
        if quizState != nil {
            if playerAnswer, exists := quizState.Answers[playerID]; exists {
                h.broadcastQuizAnswerUpdate(playerAnswer)
            }
        }
    }
}

func (h *Handler) handleQuizNext(conn *websocket.Conn) {
    question := h.server.NextQuizQuestion()
    if question != nil {
        h.broadcastQuizQuestion(question)
    }
}
```

- [ ] **Step 3: Implement broadcast functions**

Add to `backend/websocket/handler.go`:

```go
func (h *Handler) broadcastQuizQuestion(question *quiz.Question) {
    payload := map[string]any{
        "id":       question.ID,
        "type":     string(question.Type),
        "question": question.Question,
        "options":  question.Options,
    }

    msg := Message{
        Type:    "quiz_question",
        Payload: payload,
    }

    h.broadcast(msg)
}

func (h *Handler) broadcastQuizAnswerUpdate(playerAnswer *game.PlayerAnswer) {
    // Find player
    var playerName string
    var playerColor string
    for _, p := range h.server.GetPlayers() {
        if p.ID == playerAnswer.PlayerID {
            playerName = p.Name
            playerColor = p.Color
            break
        }
    }

    payload := map[string]any{
        "playerId":   playerAnswer.PlayerID,
        "playerName": playerName,
        "playerColor": playerColor,
        "answer":     playerAnswer.Answer,
        "status":     string(playerAnswer.Status),
        "isWinner":   playerAnswer.IsWinner,
    }

    msg := Message{
        Type:    "quiz_answer_update",
        Payload: payload,
    }

    h.broadcast(msg)
}
```

Add import if needed:
```go
import "github.com/yockii/lan_qr/backend/quiz"
```

- [ ] **Step 4: Build and test**

Run: `cd backend && go build`
Expected: Success

- [ ] **Step 5: Commit**

```bash
git add backend/websocket/handler.go
git commit -m "feat(backend): add quiz WebSocket message handlers"
```

---

## Phase 2: Frontend - Type Definitions

### Task 7: Extend Frontend Types

**Files:**
- Modify: `frontend/src/shared/types.ts`

- [ ] **Step 1: Add quiz-related types**

Add to `frontend/src/shared/types.ts`:

```typescript
export type QuestionType = 'single_choice' | 'true_false' | 'open_ended'

export interface Question {
  id: string
  type: QuestionType
  question: string
  options: string[]
  // correctAnswer is not sent to clients
}

export interface QuizAnswerUpdate {
  playerId: string
  playerName: string
  playerColor: string
  answer: string
  status: 'pending' | 'correct' | 'incorrect'
  isWinner: boolean
}

export interface QuizQuestionMessage {
  id: string
  type: QuestionType
  question: string
  options: string[]
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/shared/types.ts
git commit -m "feat(frontend): add quiz type definitions"
```

---

## Phase 3: Frontend - Host Components

### Task 8: Create Quiz Display Component

**Files:**
- Create: `frontend/src/host/components/QuizDisplay.vue`

- [ ] **Step 1: Create QuizDisplay component**

```vue
<template>
  <div class="quiz-display bg-slate-800 rounded-2xl p-8 mb-6">
    <div v-if="!currentQuestion" class="text-center py-12">
      <button
        @click="handleStart"
        class="px-8 py-4 bg-blue-500 hover:bg-blue-600 rounded-lg text-xl font-semibold transition-colors"
      >
        开始答题
      </button>
    </div>

    <div v-else>
      <div class="text-2xl font-bold mb-8 text-center">
        {{ currentQuestion.question }}
      </div>

      <!-- Options for single choice and true/false -->
      <div
        v-if="currentQuestion.options.length > 0"
        class="grid gap-4 max-w-2xl mx-auto"
        :class="currentQuestion.options.length === 2 ? 'grid-cols-2' : 'grid-cols-1'"
      >
        <button
          v-for="(option, index) in currentQuestion.options"
          :key="index"
          class="px-6 py-4 bg-slate-700 hover:bg-slate-600 rounded-lg text-lg text-left transition-colors"
        >
          {{ option }}
        </button>
      </div>

      <!-- No options for open-ended (just show question) -->
      <div v-else class="text-center text-slate-400">
        (问答题，选手输入答案)
      </div>

      <div class="mt-8 text-center">
        <button
          @click="handleNext"
          class="px-8 py-4 bg-green-500 hover:bg-green-600 rounded-lg text-xl font-semibold transition-colors"
        >
          下一题
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Question } from '@/shared/types'

const props = defineProps<{
  ws: WebSocketClient | null
}>()

const currentQuestion = ref<Question | null>(null)

const handleStart = () => {
  if (props.ws) {
    props.ws.send('quiz_start', {})
  }
}

const handleNext = () => {
  if (props.ws) {
    props.ws.send('quiz_next', {})
  }
}

// Listen for quiz questions
if (props.ws) {
  props.ws.on('quiz_question', (payload: Question) => {
    currentQuestion.value = payload
  })
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/host/components/QuizDisplay.vue
git commit -m "feat(frontend): add QuizDisplay component for host"
```

### Task 9: Modify PlayerList for Quiz Mode

**Files:**
- Modify: `frontend/src/host/components/PlayerList.vue`

- [ ] **Step 1: Read existing PlayerList component**

Read the file to understand current structure

- [ ] **Step 2: Add quiz answer display**

Modify the player display to show answers:

```vue
<template>
  <div class="bg-slate-800 rounded-2xl p-6">
    <h3 class="text-xl font-bold mb-4">选手列表 ({{ players.length }})</h3>

    <div class="space-y-2">
      <div
        v-for="player in players"
        :key="player.id"
        class="flex items-center justify-between p-3 bg-slate-700 rounded-lg"
      >
        <div class="flex items-center gap-3">
          <span v-if="player.isWinner" class="text-2xl">👑</span>
          <span class="font-semibold" :style="{ color: player.color }">
            {{ player.name }}
          </span>
          <span class="text-xs text-slate-400">({{ player.deviceType === 'mobile' ? '手机' : '电脑' }})</span>
        </div>

        <div class="flex items-center gap-2">
          <!-- Quiz mode: show answer -->
          <template v-if="player.answer">
            <span
              v-if="player.answerStatus === 'pending'"
              @click="handleJudge(player.id, 'correct')"
              class="px-3 py-1 bg-yellow-500 hover:bg-yellow-600 rounded cursor-pointer transition-colors"
            >
              {{ player.answer }} (待确认)
            </span>
            <span
              v-else-if="player.answerStatus === 'correct'"
              @click="handleJudge(player.id, 'incorrect')"
              class="px-3 py-1 bg-green-500 hover:bg-green-600 rounded cursor-pointer transition-colors"
            >
              ✓{{ player.answer }}
            </span>
            <span
              v-else
              @click="handleJudge(player.id, 'pending')"
              class="px-3 py-1 bg-red-500 hover:bg-red-600 rounded cursor-pointer transition-colors"
            >
              ✗{{ player.answer }}
            </span>
          </template>

          <!-- Buzzer mode: remove button -->
          <button
            v-if="!isQuizMode"
            @click="$emit('remove-player', player.id)"
            class="px-3 py-1 bg-red-500 hover:bg-red-600 rounded text-sm transition-colors"
          >
            移除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Player } from '@/shared/types'

const props = defineProps<{
  players: Player[]
  isQuizMode: boolean
  ws: any
}>()

defineEmits<{
  removePlayer: [playerId: string]
}>()

const handleJudge = (playerId: string, status: string) => {
  if (props.ws) {
    props.ws.send('quiz_judge', {
      playerId: playerId,
      correct: status === 'correct'
    })
  }
}
</script>
```

Note: This is a simplified version - adjust based on actual PlayerList structure

- [ ] **Step 3: Commit**

```bash
git add frontend/src/host/components/PlayerList.vue
git commit -m "feat(frontend): add quiz answer display to PlayerList"
```

### Task 10: Modify HostApp for Quiz Mode

**Files:**
- Modify: `frontend/src/host/HostApp.vue`

- [ ] **Step 1: Read existing HostApp**

Read the file to understand current structure

- [ ] **Step 2: Add quiz mode detection and QuizDisplay component**

Add quiz mode detection and component:

```typescript
import QuizDisplay from './components/QuizDisplay.vue'

const gameMode = ref<'buzzer' | 'quiz'>('buzzer')

// Listen for game mode
ws.on('quiz_question', () => {
  gameMode.value = 'quiz'
})
```

Update template to conditionally render:
```vue
<QuizDisplay v-if="gameMode === 'quiz'" :ws="ws" />

<BuzzerDisplay
  v-else
  :game-state="gameState"
  :winner="winner"
/>

<PlayerList
  :players="enhancedPlayers"
  :is-quiz-mode="gameMode === 'quiz'"
  :ws="ws"
  @remove-player="handleRemovePlayer"
/>
```

Note: Need to create enhancedPlayers that includes quiz answers - adjust based on actual implementation

- [ ] **Step 3: Commit**

```bash
git add frontend/src/host/HostApp.vue
git commit -m "feat(frontend): add quiz mode support to HostApp"
```

---

## Phase 4: Frontend - Player Components

### Task 11: Create Quiz Answer Component

**Files:**
- Create: `frontend/src/player/components/QuizAnswer.vue`

- [ ] **Step 1: Create QuizAnswer component**

```vue
<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <div class="text-center py-6 bg-slate-800">
      <div class="text-2xl font-bold">{{ playerName }}</div>
      <div class="mt-2">
        <span
          class="inline-block px-4 py-1 rounded-full text-sm"
          :class="connected ? 'bg-green-500' : 'bg-red-500'"
        >
          {{ connected ? '✓ 已连接' : '✗ 未连接' }}
        </span>
      </div>
    </div>

    <!-- Waiting for question -->
    <div v-if="!currentQuestion" class="flex-1 flex items-center justify-center px-8">
      <div class="text-center text-slate-400">
        <div class="text-4xl mb-4">⏳</div>
        <div class="text-xl">等待主持人开始...</div>
      </div>
    </div>

    <!-- Answer submitted -->
    <div v-else-if="answerSubmitted" class="flex-1 flex items-center justify-center px-8">
      <div class="text-center text-slate-400">
        <div class="text-4xl mb-4">✓</div>
        <div class="text-xl">等待下一题</div>
      </div>
    </div>

    <!-- Show question and input -->
    <div v-else class="flex-1 flex flex-col items-center justify-center px-8 py-8">
      <div class="text-3xl font-bold mb-8 text-center">
        {{ currentQuestion.question }}
      </div>

      <!-- Multiple choice / true false -->
      <div
        v-if="currentQuestion.options.length > 0"
        class="w-full max-w-md grid gap-4"
        :class="currentQuestion.options.length === 2 ? 'grid-cols-2' : 'grid-cols-1'"
      >
        <button
          v-for="(option, index) in currentQuestion.options"
          :key="index"
          @click="handleAnswer(option)"
          class="w-full px-6 py-4 bg-slate-800 hover:bg-slate-700 rounded-lg text-lg text-left transition-colors"
        >
          {{ option }}
        </button>
      </div>

      <!-- Open-ended -->
      <div v-else class="w-full max-w-md">
        <input
          v-model="openAnswer"
          @keyup.enter="handleAnswer(openAnswer)"
          type="text"
          placeholder="输入你的答案..."
          class="w-full px-6 py-4 text-xl bg-slate-800 border-2 border-slate-700 rounded-lg text-white text-center focus:outline-none focus:border-blue-500"
        />
        <button
          @click="handleAnswer(openAnswer)"
          :disabled="!openAnswer.trim()"
          class="w-full mt-4 px-8 py-4 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-700 disabled:cursor-not-allowed rounded-lg text-xl font-semibold transition-colors"
        >
          提交答案
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Question } from '@/shared/types'

const props = defineProps<{
  playerName: string
  ws: WebSocketClient
  connected: boolean
}>()

const currentQuestion = ref<Question | null>(null)
const answerSubmitted = ref(false)
const openAnswer = ref('')

const handleAnswer = (answer: string) => {
  if (!answer.trim()) return

  props.ws.send('quiz_answer', { answer })
  answerSubmitted.value = true
  openAnswer.value = ''
}

// Listen for quiz questions
props.ws.on('quiz_question', (payload: Question) => {
  currentQuestion.value = payload
  answerSubmitted.value = false
})

// Listen for next question (reset state)
props.ws.on('quiz_next', () => {
  currentQuestion.value = null
  answerSubmitted.value = false
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/player/components/QuizAnswer.vue
git commit -m "feat(frontend): add QuizAnswer component for players"
```

### Task 12: Modify PlayerApp for Quiz Mode

**Files:**
- Modify: `frontend/src/player/PlayerApp.vue`

- [ ] **Step 1: Read existing PlayerApp**

Read the file to understand current structure

- [ ] **Step 2: Add quiz mode detection**

Add quiz mode and conditionally render:

```typescript
import QuizAnswer from './components/QuizAnswer.vue'

const gameMode = ref<'buzzer' | 'quiz'>('buzzer')
const currentStep = ref<'input' | 'game'>('input')

// Listen for game mode
ws.on('quiz_question', () => {
  gameMode.value = 'quiz'
  currentStep.value = 'game'
})
```

Update template:
```vue
<NameInput v-if="currentStep === 'input'" />
<BuzzerButton v-else-if="gameMode === 'buzzer'" :player-name="playerName" :ws="ws" :connected="connected" />
<QuizAnswer v-else :player-name="playerName" :ws="ws" :connected="connected" />
```

Note: Adjust based on actual BuzzerButton props

- [ ] **Step 3: Commit**

```bash
git add frontend/src/player/PlayerApp.vue
git commit -m "feat(frontend): add quiz mode support to PlayerApp"
```

---

## Phase 5: Integration and Testing

### Task 13: Create Sample Question File

**Files:**
- Create: `questions.txt`

- [ ] **Step 1: Create sample questions**

Create `questions.txt` in project root:

```
[单选]
什么是CPU?|A.中央处理器|B.显卡|C.内存|D.硬盘|A

[单选]
下列哪个是编程语言？|A.HTML|B.Python|C.JPG|D.MP3|B

[判断]
地球是圆的|对

[判断]
太阳从西边升起|错

[问答]
中国的首都在哪里？|北京

[问答]
一年有多少个月？|12个
```

- [ ] **Step 2: Commit**

```bash
git add questions.txt
git commit -m "feat: add sample question bank"
```

### Task 14: Build and Test

**Files:**
- Build: `build/lan-buzzer.exe`

- [ ] **Step 1: Build frontend**

Run: `cd frontend && npx vite build`
Expected: Success

- [ ] **Step 2: Copy to backend**

Run: `rm -rf backend/embed/dist/* && cp -r frontend/dist/. backend/embed/dist/`
Expected: Success

- [ ] **Step 3: Build backend**

Run: `cd backend && go build -o ../build/lan-buzzer.exe .`
Expected: Success

- [ ] **Step 4: Test buzzer mode (no questions.txt)**

1. Rename questions.txt temporarily:
   ```bash
   mv questions.txt questions.txt.bak
   ```

2. Run build/lan-buzzer.exe

3. Verify buzzer mode works (existing functionality)

4. Rename back:
   ```bash
   mv questions.txt.bak questions.txt
   ```

- [ ] **Step 5: Test quiz mode**

1. Run build/lan-buzzer.exe

2. Verify:
   - Browser opens to quiz mode
   - Host sees "开始答题" button
   - Player joins and sees question interface
   - Host clicks "开始答题"
   - Players see question and options
   - Players submit answers
   - Host sees answers appear
   - For open-ended, host can click to judge
   - First correct answer shows 👑
   - Host clicks "下一题"
   - New question appears (no repeats in round)

- [ ] **Step 6: Test question rotation**

1. Answer all questions
2. Verify questions don't repeat in first round
3. After all questions answered, verify it starts over

- [ ] **Step 7: Fix any issues found**

Document any bugs found and fixes applied

- [ ] **Step 8: Commit build**

```bash
git add -A
git commit -m "feat: complete quiz mode implementation"
```

---

## Phase 6: Polish and Documentation

### Task 15: Update Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add quiz mode documentation**

Add to README.md:

```markdown
## Quiz Mode

The system supports two modes:

### Buzzer Mode (Default)
Traditional buzzer system where players race to buzz in first.

### Quiz Mode
When a `questions.txt` file is present, the system automatically enters quiz mode:

1. Create a `questions.txt` file in the same directory as `lan-buzzer.exe`
2. Format:
   ```
   [单选]
   问题|选项A|选项B|选项C|选项D|正确答案

   [判断]
   问题|正确答案

   [问答]
   问题|正确答案
   ```
3. Run `lan-buzzer.exe`
4. Host clicks "开始答题" to start
5. Players see questions on their devices
6. Multiple choice/true/false: auto-judged
7. Open-ended: host clicks answer to judge (待确认 → ✓ → ✗)
8. First correct answer wins (👑)
9. Host clicks "下一题" for next question
```

- [ ] **Step 2: Update CLAUDE.md**

Add quiz mode information to architecture section

- [ ] **Step 3: Commit**

```bash
git add README.md CLAUDE.md
git commit -m "docs: add quiz mode documentation"
```

### Task 16: Final Testing and Release

**Files:**
- Release: `build/lan-buzzer.exe`

- [ ] **Step 1: Complete end-to-end test**

Test all features with multiple devices:
- Buzzer mode (without questions.txt)
- Quiz mode with all question types
- Question rotation
- Manual judging
- Mode switching

- [ ] **Step 2: Create release build**

Run: `./build.sh`

- [ ] **Step 3: Tag release**

```bash
git tag -a v1.1.0 -m "Add quiz mode support"
git push origin feature/question-bank --tags
```

- [ ] **Step 4: Merge to main**

```bash
git checkout main
git merge feature/question-bank
git push origin main
```

- [ ] **Step 5: Create GitHub release**

Use gh CLI or GitHub web interface to create release v1.1.0

---

## Summary

This plan implements quiz mode for the LAN Buzzer system with:

**Backend:**
- Question data structures and parser for txt format
- Question bank with deduplication and random selection
- Game server extensions for quiz state management
- WebSocket message handlers for quiz flow

**Frontend:**
- Type definitions for quiz-related data
- QuizDisplay component for host
- Enhanced PlayerList with answer display and manual judging
- QuizAnswer component for players
- Mode detection in HostApp and PlayerApp

**Features:**
- Automatic mode detection (quiz vs buzzer)
- Support for single choice, true/false, and open-ended questions
- Auto-judging for choice/TF, manual judging for open-ended
- Real-time answer display
- Winner identification (👑)
- Question rotation with no repeats in round

**Testing:**
- Unit tests for parser and question bank
- Integration tests for quiz flow
- Manual E2E testing with multiple devices
- Mode switching verification
