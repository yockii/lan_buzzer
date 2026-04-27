# Question Bank Feature Design

## Overview

Add a quiz mode to the LAN Buzzer system where players compete to answer questions correctly and quickly. The system automatically loads questions from a `questions.txt` file and presents them one by one. Players submit answers through their devices, and the host sees real-time results.

## Goals

1. Enable quiz-based competition alongside the existing buzzer mode
2. Support multiple question types: multiple choice, true/false, and open-ended
3. Automatic question randomization with deduplication within each round
4. Real-time answer display on host screen
5. Manual judging for open-ended questions

## Architecture

### Backend Structure

```
backend/
├── quiz/
│   ├── question.go       # Question data structures
│   ├── question_bank.go  # Question bank management (loading, deduplication, random selection)
│   └── parser.go         # txt file parsing
├── game/
│   └── server.go         # Extended to support quiz mode
└── websocket/
    └── handler.go        # Extended with quiz message types
```

### Frontend Structure

```
frontend/src/
├── host/
│   ├── HostApp.vue              # Modified to detect mode
│   └── components/
│       ├── QuizDisplay.vue      # NEW: Question display for host
│       └── PlayerList.vue       # Modified to show answers
└── player/
    ├── PlayerApp.vue            # Modified to detect mode
    └── components/
        ├── QuizAnswer.vue       # NEW: Answer submission interface
        └── BuzzerButton.vue     # Existing: buzzer mode
```

## Components

### 1. Question Bank (backend/quiz/question_bank.go)

**Responsibilities:**
- Load questions from `questions.txt` on startup
- Parse question file format
- Maintain list of asked question IDs for deduplication
- Provide random unasked question
- Reset asked list when all questions exhausted

**Interface:**
```go
type QuestionBank struct {
    questions []Question
    askedIDs  map[string]bool
}

func NewQuestionBank(filePath string) (*QuestionBank, error)
func (qb *QuestionBank) GetRandomQuestion() *Question
func (qb *QuestionBank) MarkAsked(questionID string)
func (qb *QuestionBank) ResetAsked()
func (qb *QuestionBank) HasQuestions() bool
```

### 2. Question Parser (backend/quiz/parser.go)

**Responsibilities:**
- Parse txt file format
- Validate question structure
- Return structured Question objects

**File Format:**
```
[单选]
什么是CPU?|A.中央处理器|B.显卡|C.内存|D.硬盘|A

[判断]
地球是圆的|对

[问答]
中国的首都在哪里？|北京
```

**Parsing Rules:**
- Lines starting with `[题型]` define question type
- Empty lines are skipped
- Fields separated by `|`
- Question types: 单选, 判断, 问答

### 3. Quiz Display (frontend/src/host/components/QuizDisplay.vue)

**Responsibilities:**
- Display current question and options
- Show "开始答题" button before question starts
- Show "下一题" button after answering begins

**Layout:**
```
┌─────────────────────────────────┐
│                                 │
│  什么是CPU?                     │
│                                 │
│  ⚪ A. 中央处理器               │
│  ⚪ B. 显卡                     │
│  ⚪ C. 内存                     │
│  ⚪ D. 硬盘                     │
│                                 │
│  [开始答题] / [下一题]          │
└─────────────────────────────────┘
```

### 4. Player List Enhancement (frontend/src/host/components/PlayerList.vue)

**New Display Format:**
- Player name + color
- Answer status: 未答 / A / B / 对 / 北京 / 待确认 / ✓正确 / ✗错误
- Click answer to cycle: 待确认 → ✓ → ✗ → 待确认
- First correct answer: 👑 icon before name

**Example:**
```
👑 张三 (红色) - ✓A
李四 (蓝色) - B
王五 (绿色) - 待确认
```

### 5. Quiz Answer Interface (frontend/src/player/components/QuizAnswer.vue)

**Responsibilities:**
- Display question and options/input
- Submit answer to server
- Show "等待下一题" after submission

**Multiple Choice UI:**
```
┌─────────────────────────┐
│  什么是CPU?             │
│                         │
│  [ A. 中央处理器 ]      │
│  [ B. 显卡 ]            │
│  [ C. 内存 ]            │
│  [ D. 硬盘 ]            │
└─────────────────────────┘
```

**True/False UI:**
```
┌─────────────────────────┐
│  地球是圆的?            │
│                         │
│  [ 对 ]  [ 错 ]         │
└─────────────────────────┘
```

**Open-ended UI:**
```
┌─────────────────────────┐
│  中国的首都在哪里？     │
│                         │
│  [  答案输入框    ]     │
│  [     提交     ]       │
└─────────────────────────┘
```

## Data Flow

### Game Flow

1. **Startup:**
   - Backend loads `questions.txt`
   - If file exists and valid → quiz mode
   - If file missing or invalid → buzzer mode (existing)

2. **Start Question:**
   - Host clicks "开始答题"
   - Backend selects random unasked question
   - Broadcasts `quiz_question` message to all clients
   - Frontend displays question

3. **Answer Submission:**
   - Player selects answer or inputs text
   - Client sends `quiz_answer` message
   - Backend records answer with timestamp
   - Broadcasts `quiz_answer_update` to host
   - Host sees answer appear next to player name

4. **Automatic Judging (Multiple Choice/True/False):**
   - Backend compares answer to correct answer
   - If correct and first → mark as winner (👑)
   - Broadcasts updated status

5. **Manual Judging (Open-ended):**
   - Answer doesn't match exactly
   - Status shows "待确认"
   - Host clicks answer to cycle: 待确认 → ✓ → ✗
   - Client sends `quiz_judge` message
   - Backend updates status
   - If marked as correct and first → mark as winner (👑)

6. **Next Question:**
   - Host clicks "下一题"
   - Backend marks current question as asked
   - Selects new random unasked question
   - If all questions asked → reset asked list
   - Broadcasts new question

### WebSocket Messages

**Client → Server:**
- `quiz_answer` - Player submits answer
  ```json
  {
    "type": "quiz_answer",
    "payload": {
      "answer": "A"
    }
  }
  ```

- `quiz_judge` - Host manually judges open-ended answer
  ```json
  {
    "type": "quiz_judge",
    "payload": {
      "playerId": "xxx",
      "correct": true
    }
  }
  ```

**Server → Client:**
- `quiz_question` - Broadcast new question
  ```json
  {
    "type": "quiz_question",
    "payload": {
      "id": "q1",
      "type": "single_choice",
      "question": "什么是CPU?",
      "options": ["A.中央处理器", "B.显卡", "C.内存", "D.硬盘"],
      "correctAnswer": "A"
    }
  }
  ```

- `quiz_answer_update` - Broadcast answer update
  ```json
  {
    "type": "quiz_answer_update",
    "payload": {
      "playerId": "xxx",
      "answer": "A",
      "status": "correct",
      "isWinner": true
    }
  }
  ```

## Error Handling

### Question File Errors
- File not found → Log warning, use buzzer mode
- Parse error → Log error details, use buzzer mode
- Empty file → Log warning, use buzzer mode
- Invalid format → Log line number and error, skip question

### Runtime Errors
- No questions available → Log error, show error message to host
- All questions asked → Automatically reset asked list, start over
- Invalid answer format → Log warning, don't crash

## State Management

### Quiz State (backend/game/server.go)
```go
type QuizState struct {
    CurrentQuestion *Question
    Answers         map[string]*PlayerAnswer  // playerId -> answer
    WinnerID        string  // First correct answer
    AskedIDs        []string
}

type PlayerAnswer struct {
    PlayerID   string
    Answer     string
    Timestamp  int64
    Status     string  // "pending", "correct", "incorrect"
    IsWinner   bool
}
```

### Game Mode Detection
```go
func (gs *GameServer) DetectMode() string {
    if gs.QuestionBank.HasQuestions() {
        return "quiz"
    }
    return "buzzer"
}
```

## Testing

### Unit Tests
- Question parser with valid/invalid formats
- Question bank random selection and deduplication
- Answer comparison logic (exact match for choice/TF)

### Integration Tests
- Load question file and enter quiz mode
- Player answer submission flow
- Host manual judging flow
- Question rotation and reset

### Manual Testing
- Create sample `questions.txt` with all question types
- Test with multiple devices
- Verify answer display and judging
- Test question rotation (no repeats in round)

## Migration Path

### Phase 1: Backend
1. Create quiz package with question structures
2. Implement parser for txt format
3. Integrate question bank into game server
4. Add quiz WebSocket messages

### Phase 2: Host Frontend
1. Create QuizDisplay component
2. Modify PlayerList to show answers
3. Add mode detection (quiz vs buzzer)

### Phase 3: Player Frontend
1. Create QuizAnswer component
2. Implement answer submission
3. Add mode detection

### Phase 4: Integration
1. End-to-end testing
2. Polish UI/UX
3. Add error handling

## Future Enhancements (Out of Scope)

- Question categories/tags
- Difficulty levels
- Scoring/points system
- Timer per question
- Question statistics
- Export/import question banks
- Rich text questions (images, audio, video)
