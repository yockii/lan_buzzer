# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LAN Buzzer is a real-time buzzer system for local networks. It supports 5-10 players with WebSocket communication for instant responses.

## Common Commands

### Development

```bash
# Backend development
cd backend
go run main.go

# Frontend development
cd frontend
npm run dev

# Build everything
./build.sh

# Run tests
cd backend && go test ./...
```

### Building

```bash
# Build frontend and backend into single binary
./build.sh

# Output: build/lan-buzzer.exe (Windows) or build/lan-buzzer (Linux/Mac)
```

## Architecture

### Game Modes

The system supports two game modes:

1. **Buzzer Mode**: Players race to buzz in first. WebSocket events: `buzz`, `buzz_result`, `state_changed`
2. **Quiz Mode**: Players answer questions from a question bank. WebSocket events: `quiz_start`, `quiz_question`, `quiz_answer`, `quiz_judge`, `quiz_next`, `quiz_answer_update`

Mode is determined by presence of `questions.txt` file in the backend directory.

### Quiz Mode Components

- **Backend**:
  - `backend/quiz/`: Question parsing and bank management
  - `backend/game/types.go`: QuizState, PlayerAnswer, AnswerStatus types
  - `backend/game/server.go`: Quiz mode logic in GameServer

- **Frontend Host**:
  - `frontend/src/host/components/QuizDisplay.vue`: Shows questions to host
  - `frontend/src/host/components/PlayerList.vue`: Shows player answers with manual judging

- **Frontend Player**:
  - `frontend/src/player/components/QuizAnswer.vue`: Player interface for answering questions

### Question Format

Questions are stored in `questions.txt` with pipe-separated values:
- `[单选]`: Single choice - question|optionA|optionB|optionC|optionD|answer
- `[判断]`: True/false - question|answer
- `[问答]`: Open-ended - question|answer

### Backend (Go + Fiber v2)

- `backend/main.go` - Entry point, serves frontend, auto-opens browser
- `backend/game/` - Core game logic and state management
- `backend/websocket/` - WebSocket handler for real-time communication
- `backend/embed/dist/` - Frontend build output (embedded in binary)

### Frontend (Vue 3 + TypeScript)

- `frontend/src/host/` - Host control panel
- `frontend/src/player/` - Player buzzer interface
- `frontend/src/shared/` - Shared types and WebSocket client

## Key Design Decisions

1. **Single Binary**: Frontend built into Go binary via go:embed - zero dependencies
2. **WebSocket**: Real-time bidirectional communication for instant buzz responses
3. **Color-Based IDs**: Players assigned colors on join, supports duplicate names
4. **Device Detection**: Auto-routes to host or player interface based on User-Agent
5. **State Machine**: Three states (waiting/ready/locked) managed by GameServer

## Game Flow

1. Host starts server, browser opens automatically
2. Players join via QR code or URL
3. Host clicks "开始抢答" → state changes to "ready"
4. Players buzz (space/enter/touch)
5. First buzz wins → state changes to "locked"
6. Host clicks "下一题" → resets to "waiting"

## Testing

- Backend: `go test ./...` in backend directory
- Frontend: Manual testing in browser (no automated tests yet)
- E2E: Run binary, connect multiple devices, test full game flow
