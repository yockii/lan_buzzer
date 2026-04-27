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
