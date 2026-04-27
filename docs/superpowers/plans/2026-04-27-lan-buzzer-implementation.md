# LAN Buzzer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a LAN-based buzzer system with Go backend and Vue3 frontend, supporting 5-10 players with real-time buzz responses.

**Architecture:** Single Go binary with embedded Vue3 static assets, WebSocket communication for real-time updates, color-based player identification to support duplicate names.

**Tech Stack:** Go 1.25+ + Fiber v3 + WebSocket, Vue 3 + TypeScript + Vite + shadcn-vue + TailwindCSS

---

## File Structure

### Backend
```
backend/
├── go.mod                      # Go module definition
├── go.sum                      # Dependency lock file
├── main.go                     # Entry point, static file serving
├── websocket/
│   └── handler.go              # WebSocket connection handler
├── game/
│   ├── types.go                # Core data structures
│   ├── server.go               # Game state management
│   └── colors.go               # Player color assignments
└── embed/
    └── static/                 # Frontend build output (embedded)
```

### Frontend
```
frontend/
├── package.json                # Node dependencies
├── vite.config.ts              # Vite configuration
├── tsconfig.json               # TypeScript config
├── tailwind.config.js          # TailwindCSS config
├── index.html                  # Entry HTML
├── src/
│   ├── main.ts                 # App bootstrap
│   ├── App.vue                 # Root component
│   ├── shared/
│   │   ├── types.ts            # Shared TypeScript types
│   │   └── websocket.ts        # WebSocket client
│   ├── host/
│   │   ├── HostApp.vue         # Host main page
│   │   └── components/
│   │       ├── ControlPanel.vue    # Start/Next buttons
│   │       ├── BuzzerDisplay.vue   # Winner display
│   │       └── PlayerList.vue      # Connected players
│   └── player/
│       ├── PlayerApp.vue       # Player main page
│       └── components/
│           ├── NameInput.vue   # Registration screen
│           ├── BuzzerButton.vue # Buzzer button
│           └── StatusModal.vue # Early buzz warning
└── dist/                       # Build output (embedded by Go)
```

---

## Task 1: Initialize Go Backend Project

**Files:**
- Create: `backend/go.mod`
- Create: `backend/go.sum`
- Create: `backend/main.go`

- [ ] **Step 1: Initialize Go module**

Run: `cd backend && go mod init github.com/yockii/lan-buzzer`

- [ ] **Step 2: Create go.mod with dependencies**

Create `backend/go.mod`:
```go
module github.com/yockii/lan-buzzer

go 1.25

require (
    github.com/gofiber/fiber/v3 v3.0.0
    github.com/gofiber/websocket/v2 v2.2.1
)
```

- [ ] **Step 3: Download dependencies**

Run: `cd backend && go mod tidy`

Expected: `go.sum` file created with dependencies

- [ ] **Step 4: Create main.go with basic Fiber server**

Create `backend/main.go`:
```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })

    app.Use(logger.New())
    app.Use(recover.New())

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("LAN Buzzer Server")
    })

    log.Fatal(app.Listen(":3000"))
}
```

- [ ] **Step 5: Test server starts**

Run: `cd backend && go run main.go`

Expected: Server starts on port 3000

Stop with: `Ctrl+C`

- [ ] **Step 6: Commit**

```bash
git add backend/go.mod backend/go.sum backend/main.go
git commit -m "feat: initialize Go backend with Fiber server"
```

---

## Task 2: Define Core Game Types

**Files:**
- Create: `backend/game/types.go`

- [ ] **Step 1: Write types test file**

Create `backend/game/types_test.go`:
```go
package game

import "testing"

func TestGameStateValues(t *testing.T) {
    if StateWaiting != "waiting" {
        t.Errorf("Expected StateWaiting to be 'waiting', got %s", StateWaiting)
    }
    if StateReady != "ready" {
        t.Errorf("Expected StateReady to be 'ready', got %s", StateReady)
    }
    if StateLocked != "locked" {
        t.Errorf("Expected StateLocked to be 'locked', got %s", StateLocked)
    }
}

func TestPlayerCreation(t *testing.T) {
    player := &Player{
        ID:         "test-id",
        Name:       "Test Player",
        Color:      "#ef4444",
        DeviceType: "desktop",
    }

    if player.ID != "test-id" {
        t.Errorf("Expected ID to be 'test-id', got %s", player.ID)
    }
    if player.Name != "Test Player" {
        t.Errorf("Expected Name to be 'Test Player', got %s", player.Name)
    }
    if player.Color != "#ef4444" {
        t.Errorf("Expected Color to be '#ef4444', got %s", player.Color)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend/game && go test -v`

Expected: FAIL with "undefined: StateWaiting"

- [ ] **Step 3: Implement types**

Create `backend/game/types.go`:
```go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend/game && go test -v`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/game/types.go backend/game/types_test.go
git commit -m "feat: define core game types"
```

---

## Task 3: Implement Player Color System

**Files:**
- Create: `backend/game/colors.go`

- [ ] **Step 1: Write colors test**

Create `backend/game/colors_test.go`:
```go
package game

import "testing"

func TestGetPlayerColor(t *testing.T) {
    colors := make(map[string]string)

    // First player gets red
    color1 := GetPlayerColor("player1", colors)
    if color1 != "#ef4444" {
        t.Errorf("Expected first color to be #ef4444, got %s", color1)
    }

    colors["player1"] = color1

    // Second player gets orange
    color2 := GetPlayerColor("player2", colors)
    if color2 != "#f97316" {
        t.Errorf("Expected second color to be #f97316, got %s", color2)
    }
}

func TestColorCycling(t *testing.T) {
    colors := make(map[string]string)

    // Add 8 players to use all colors
    for i := 0; i < 8; i++ {
        id := string(rune('a' + i))
        colors[id] = GetPlayerColor(id, colors)
    }

    // 9th player should cycle back to red
    color9 := GetPlayerColor("player9", colors)
    if color9 != "#ef4444" {
        t.Errorf("Expected 9th color to cycle to #ef4444, got %s", color9)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend/game && go test -v`

Expected: FAIL with "undefined: GetPlayerColor"

- [ ] **Step 3: Implement color system**

Create `backend/game/colors.go`:
```go
package game

var PlayerColors = []string{
    "#ef4444", // red
    "#f97316", // orange
    "#eab308", // yellow
    "#22c55e", // green
    "#06b6d4", // cyan
    "#3b82f6", // blue
    "#a855f7", // purple
    "#ec4899", // pink
}

func GetPlayerColor(playerID string, existingColors map[string]string) string {
    usedColors := make(map[string]bool)
    for _, color := range existingColors {
        usedColors[color] = true
    }

    for _, color := range PlayerColors {
        if !usedColors[color] {
            return color
        }
    }

    // All colors used, cycle back to first
    return PlayerColors[0]
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend/game && go test -v`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/game/colors.go backend/game/colors_test.go
git commit -m "feat: implement player color assignment system"
```

---

## Task 4: Implement Game Server Logic

**Files:**
- Create: `backend/game/server.go`

- [ ] **Step 1: Write server tests**

Create `backend/game/server_test.go`:
```go
package game

import (
    "testing"
    "time"
)

func TestNewGameServer(t *testing.T) {
    server := NewGameServer()

    if server.State != StateWaiting {
        t.Errorf("Expected initial state to be waiting, got %s", server.State)
    }

    if len(server.Players) != 0 {
        t.Errorf("Expected no players initially, got %d", len(server.Players))
    }

    if server.WinnerID != nil {
        t.Errorf("Expected no winner initially, got %v", server.WinnerID)
    }
}

func TestAddPlayer(t *testing.T) {
    server := NewGameServer()

    player := &Player{
        ID:         "test1",
        Name:       "Test Player",
        Color:      "#ef4444",
        DeviceType: "desktop",
    }

    server.AddPlayer(player)

    if len(server.Players) != 1 {
        t.Errorf("Expected 1 player, got %d", len(server.Players))
    }

    if server.Players["test1"] != player {
        t.Error("Player not added correctly")
    }
}

func TestStartGame(t *testing.T) {
    server := NewGameServer()

    server.StartGame()

    if server.State != StateReady {
        t.Errorf("Expected state to be ready, got %s", server.State)
    }

    if server.StartTime == nil {
        t.Error("Expected start time to be set")
    }
}

func TestRecordBuzz(t *testing.T) {
    server := NewGameServer()
    server.StartGame()

    winnerID := "player1"
    server.RecordBuzz(winnerID)

    if server.State != StateLocked {
        t.Errorf("Expected state to be locked, got %s", server.State)
    }

    if server.WinnerID == nil {
        t.Error("Expected winner to be set")
    }

    if *server.WinnerID != winnerID {
        t.Errorf("Expected winner to be %s, got %s", winnerID, *server.WinnerID)
    }
}

func TestResetGame(t *testing.T) {
    server := NewGameServer()
    server.StartGame()

    winnerID := "player1"
    server.RecordBuzz(winnerID)

    server.ResetGame()

    if server.State != StateWaiting {
        t.Errorf("Expected state to be waiting, got %s", server.State)
    }

    if server.WinnerID != nil {
        t.Error("Expected winner to be cleared")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend/game && go test -v`

Expected: FAIL with "undefined: NewGameServer"

- [ ] **Step 3: Implement game server**

Create `backend/game/server.go`:
```go
package game

import (
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd backend/game && go test -v`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/game/server.go backend/game/server_test.go
git commit -m "feat: implement game server state management"
```

---

## Task 5: Implement WebSocket Handler

**Files:**
- Create: `backend/websocket/handler.go`

- [ ] **Step 1: Create WebSocket message types**

Create `backend/websocket/types.go`:
```go
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
```

- [ ] **Step 2: Implement WebSocket handler**

Create `backend/websocket/handler.go`:
```go
package websocket

import (
    "encoding/json"
    "log"
    "sync"

    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/websocket/v2"
    "github.com/google/uuid"
    "github.com/yockii/lan-buzzer/backend/game"
)

type Handler struct {
    server  *game.GameServer
    clients map[string]*websocket.Conn
    mutex   sync.RWMutex
}

func NewHandler(gs *game.GameServer) *Handler {
    return &Handler{
        server:  gs,
        clients: make(map[string]*websocket.Conn),
    }
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
    app.Use("/ws", func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            c.Locals("allowed", true)
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    app.Get("/ws", websocket.New(h.handleConnection))
}

func (h *Handler) handleConnection(conn *websocket.Conn) {
    playerID := uuid.New().String()
    var player *game.Player

    defer func() {
        conn.Close()
        h.mutex.Lock()
        delete(h.clients, playerID)
        h.mutex.Unlock()

        if player != nil {
            h.server.RemovePlayer(playerID)
            h.broadcastPlayerList()
        }
        log.Printf("Player %s disconnected", playerID)
    }()

    log.Printf("New connection: %s", playerID)

    for {
        msg := new(Message)
        if err := conn.ReadJSON(msg); err != nil {
            log.Printf("Read error: %v", err)
            break
        }

        switch msg.Type {
        case "join":
            h.handleJoin(conn, playerID, msg)
        case "buzz":
            h.handleBuzz(conn, playerID)
        }
    }
}

func (h *Handler) handleJoin(conn *websocket.Conn, playerID string, msg *Message) {
    payload, ok := msg.Payload.(map[string]any)
    if !ok {
        h.sendError(conn, "Invalid payload")
        return
    }

    name, _ := payload["name"].(string)
    deviceType, _ := payload["deviceType"].(string)

    if name == "" {
        h.sendError(conn, "Name is required")
        return
    }

    if deviceType != "desktop" && deviceType != "mobile" {
        deviceType = "desktop"
    }

    h.mutex.Lock()
    h.clients[playerID] = conn
    h.mutex.Unlock()

    existingColors := make(map[string]string)
    for _, p := range h.server.GetPlayers() {
        existingColors[p.ID] = p.Color
    }

    player = &game.Player{
        ID:         playerID,
        Name:       name,
        Color:      game.GetPlayerColor(playerID, existingColors),
        DeviceType: deviceType,
    }

    h.server.AddPlayer(player)
    h.broadcastPlayerList()
    h.broadcastState(conn)
}

func (h *Handler) handleBuzz(conn *websocket.Conn, playerID string) {
    state := h.server.GetState()

    if state == game.StateWaiting {
        h.sendEarlyBuzzWarning(conn)
        return
    }

    if state == game.StateLocked {
        return
    }

    if h.server.RecordBuzz(playerID) {
        winner := h.server.GetWinner()
        h.broadcastBuzzResult(winner)
    }
}

func (h *Handler) broadcastState(conn *websocket.Conn) {
    state := h.server.GetState()

    payload := StateChangedPayload{
        State: string(state),
    }

    if state == game.StateReady && h.server.StartTime != nil {
        ts := h.server.StartTime.Unix()
        payload.StartTime = &ts
    }

    msg := Message{
        Type:    "state_changed",
        Payload: payload,
    }

    conn.WriteJSON(msg)
}

func (h *Handler) broadcastPlayerList() {
    players := h.server.GetPlayers()

    playerInfos := make([]PlayerInfo, len(players))
    for i, p := range players {
        playerInfos[i] = PlayerInfo{
            ID:         p.ID,
            Name:       p.Name,
            Color:      p.Color,
            DeviceType: p.DeviceType,
        }
    }

    payload := PlayerListPayload{
        Players: playerInfos,
    }

    msg := Message{
        Type:    "player_list",
        Payload: payload,
    }

    h.broadcast(msg)
}

func (h *Handler) broadcastBuzzResult(winner *game.Player) {
    var winnerInfo *PlayerInfo
    if winner != nil {
        winnerInfo = &PlayerInfo{
            ID:         winner.ID,
            Name:       winner.Name,
            Color:      winner.Color,
            DeviceType: winner.DeviceType,
        }
    }

    payload := BuzzResultPayload{
        Winner:  winnerInfo,
        IsEarly: false,
    }

    msg := Message{
        Type:    "buzz_result",
        Payload: payload,
    }

    h.broadcast(msg)
}

func (h *Handler) sendEarlyBuzzWarning(conn *websocket.Conn) {
    msg := Message{
        Type: "early_buzz_warning",
        Payload: map[string]string{
            "message": "请等待主持人开始！",
        },
    }

    conn.WriteJSON(msg)
}

func (h *Handler) sendError(conn *websocket.Conn, message string) {
    msg := Message{
        Type: "error",
        Payload: ErrorPayload{
            Message: message,
        },
    }

    conn.WriteJSON(msg)
}

func (h *Handler) broadcast(msg Message) {
    h.mutex.RLock()
    defer h.mutex.RUnlock()

    for _, conn := range h.clients {
        conn.WriteJSON(msg)
    }
}

func (h *Handler) StartGame() {
    h.server.StartGame()

    state := h.server.GetState()
    payload := StateChangedPayload{
        State: string(state),
    }

    if h.server.StartTime != nil {
        ts := h.server.StartTime.Unix()
        payload.StartTime = &ts
    }

    msg := Message{
        Type:    "state_changed",
        Payload: payload,
    }

    h.broadcast(msg)
}

func (h *Handler) ResetGame() {
    h.server.ResetGame()

    payload := StateChangedPayload{
        State: string(game.StateWaiting),
    }

    msg := Message{
        Type:    "state_changed",
        Payload: payload,
    }

    h.broadcast(msg)
}
```

- [ ] **Step 3: Update main.go to use WebSocket handler**

Modify `backend/main.go`:
```go
package main

import (
    "log"

    "github.com/yockii/lan-buzzer/backend/game"
    "github.com/yockii/lan-buzzer/backend/websocket"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
    "github.com/gofiber/websocket/v2"
)

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })

    app.Use(logger.New())
    app.Use(recover.New())

    gameServer := game.NewGameServer()
    wsHandler := websocket.NewHandler(gameServer)
    wsHandler.RegisterRoutes(app)

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("LAN Buzzer Server")
    })

    log.Fatal(app.Listen(":3000"))
}
```

- [ ] **Step 4: Add missing import**

Run: `cd backend && go mod tidy`

- [ ] **Step 5: Test server compiles**

Run: `cd backend && go build -o ../build/lan-buzzer.exe .`

Expected: Binary created at `build/lan-buzzer.exe`

- [ ] **Step 6: Commit**

```bash
git add backend/websocket/ backend/main.go
git commit -m "feat: implement WebSocket handler"
```

---

## Task 6: Initialize Vue3 Frontend Project

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tailwind.config.js`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`

- [ ] **Step 1: Create package.json**

Create `frontend/package.json`:
```json
{
  "name": "lan-buzzer-frontend",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.2.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "vue-tsc": "^1.8.0",
    "tailwindcss": "^3.4.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0"
  }
}
```

- [ ] **Step 2: Install dependencies**

Run: `cd frontend && npm install`

Expected: `node_modules` directory created

- [ ] **Step 3: Create vite.config.ts**

Create `frontend/vite.config.ts`:
```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true
  }
})
```

- [ ] **Step 4: Create tsconfig.json**

Create `frontend/tsconfig.json`:
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "module": "ESNext",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "preserve",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

- [ ] **Step 5: Create tsconfig.node.json**

Create `frontend/tsconfig.node.json`:
```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true
  },
  "include": ["vite.config.ts"]
}
```

- [ ] **Step 6: Create tailwind.config.js**

Create `frontend/tailwind.config.js`:
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

- [ ] **Step 7: Create postcss.config.js**

Create `frontend/postcss.config.js`:
```javascript
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

- [ ] **Step 8: Create index.html**

Create `frontend/index.html`:
```html
<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8">
    <link rel="icon" type="image/svg+xml" href="/vite.svg">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LAN Buzzer</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

- [ ] **Step 9: Create main.ts**

Create `frontend/src/main.ts`:
```typescript
import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

createApp(App).mount('#app')
```

- [ ] **Step 10: Create style.css**

Create `frontend/src/style.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
}

#app {
  width: 100vw;
  height: 100vh;
}
```

- [ ] **Step 11: Create basic App.vue**

Create `frontend/src/App.vue`:
```vue
<template>
  <div class="w-full h-full bg-slate-900 text-white">
    <h1 class="text-4xl text-center pt-20">LAN Buzzer</h1>
  </div>
</template>

<script setup lang="ts">
</script>
```

- [ ] **Step 12: Test dev server**

Run: `cd frontend && npm run dev`

Expected: Dev server starts on http://localhost:5173

Stop with: `Ctrl+C`

- [ ] **Step 13: Commit**

```bash
git add frontend/
git commit -m "feat: initialize Vue3 frontend with Vite and TailwindCSS"
```

---

## Task 7: Implement Shared Types and WebSocket Client

**Files:**
- Create: `frontend/src/shared/types.ts`
- Create: `frontend/src/shared/websocket.ts`

- [ ] **Step 1: Create shared types**

Create `frontend/src/shared/types.ts`:
```typescript
export type GameState = 'waiting' | 'ready' | 'locked'

export interface Player {
  id: string
  name: string
  color: string
  deviceType: 'desktop' | 'mobile'
  connectedAt: string
}

export interface Message {
  type: string
  payload: any
}

export interface StateChangedPayload {
  state: GameState
  startTime?: number
}

export interface PlayerListPayload {
  players: Player[]
}

export interface BuzzResultPayload {
  winner?: Player
  isEarly: boolean
}

export interface ErrorPayload {
  message: string
}
```

- [ ] **Step 2: Create WebSocket client**

Create `frontend/src/shared/websocket.ts`:
```typescript
import { ref, Ref } from 'vue'
import type { Message } from './types'

export class WebSocketClient {
  private ws: WebSocket | null = null
  private reconnectTimer: number | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5

  public connected: Ref<boolean> = ref(false)
  public messageHandlers: Map<string, (payload: any) => void> = new Map()

  constructor(private url: string) {}

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.connected.value = true
      this.reconnectAttempts = 0
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      this.connected.value = false
      this.scheduleReconnect()
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    this.ws.onmessage = (event) => {
      try {
        const message: Message = JSON.parse(event.data)
        const handler = this.messageHandlers.get(message.type)
        if (handler) {
          handler(message.payload)
        }
      } catch (error) {
        console.error('Failed to parse message:', error)
      }
    }
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.connected.value = false
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnect attempts reached')
      return
    }

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000)
    this.reconnectAttempts++

    this.reconnectTimer = window.setTimeout(() => {
      console.log(`Reconnecting... (attempt ${this.reconnectAttempts})`)
      this.connect()
    }, delay)
  }

  send(type: string, payload: any) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.error('WebSocket not connected')
      return
    }

    const message: Message = { type, payload }
    this.ws.send(JSON.stringify(message))
  }

  on(type: string, handler: (payload: any) => void) {
    this.messageHandlers.set(type, handler)
  }

  off(type: string) {
    this.messageHandlers.delete(type)
  }
}
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/shared/
git commit -m "feat: implement shared types and WebSocket client"
```

---

## Task 8: Implement Host Control Panel

**Files:**
- Create: `frontend/src/host/components/ControlPanel.vue`

- [ ] **Step 1: Create ControlPanel component**

Create `frontend/src/host/components/ControlPanel.vue`:
```vue
<template>
  <div class="flex gap-4 justify-center mb-6">
    <button
      @click="$emit('start')"
      :disabled="gameState !== 'waiting'"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
      :class="gameState === 'waiting' ? 'bg-blue-500 hover:bg-blue-600' : 'bg-gray-600 cursor-not-allowed opacity-50'"
    >
      开始抢答
    </button>

    <button
      @click="$emit('reset')"
      :disabled="gameState === 'waiting'"
      class="px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
      :class="gameState !== 'waiting' ? 'bg-gray-500 hover:bg-gray-600' : 'bg-gray-700 cursor-not-allowed opacity-50'"
    >
      下一题
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  gameState: 'waiting' | 'ready' | 'locked'
}>()

defineEmits<{
  start: []
  reset: []
}>()
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/host/components/ControlPanel.vue
git commit -m "feat: implement host control panel"
```

---

## Task 9: Implement Host Buzzer Display

**Files:**
- Create: `frontend/src/host/components/BuzzerDisplay.vue`

- [ ] **Step 1: Create BuzzerDisplay component**

Create `frontend/src/host/components/BuzzerDisplay.vue`:
```vue
<template>
  <div class="flex-1 flex items-center justify-center px-8">
    <div class="w-full text-center py-16 bg-slate-800 rounded-2xl">
      <div v-if="winner" class="space-y-4">
        <div class="text-7xl font-bold" :style="{ color: winner.color }">
          {{ winner.name }}
        </div>
        <div class="text-2xl text-slate-300">
          🎉 抢到了！
        </div>
      </div>
      <div v-else class="space-y-4">
        <div class="text-5xl font-bold text-slate-400">
          {{ stateText }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Player } from '@/shared/types'

const props = defineProps<{
  gameState: 'waiting' | 'ready' | 'locked'
  winner: Player | null
}>()

const stateText = computed(() => {
  switch (props.gameState) {
    case 'waiting':
      return '等待开始'
    case 'ready':
      return '准备抢答...'
    case 'locked':
      return '已锁定'
  }
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/host/components/BuzzerDisplay.vue
git commit -m "feat: implement host buzzer display"
```

---

## Task 10: Implement Host Player List

**Files:**
- Create: `frontend/src/host/components/PlayerList.vue`

- [ ] **Step 1: Create PlayerList component**

Create `frontend/src/host/components/PlayerList.vue`:
```vue
<template>
  <div class="text-center py-4 text-slate-400">
    <span class="text-sm">已连接 {{ players.length }} 位选手</span>
    <span v-if="players.length > 0" class="ml-4">
      <span
        v-for="(player, index) in players"
        :key="player.id"
        class="inline-flex items-center mx-1"
      >
        <span
          class="w-3 h-3 rounded-full mr-1"
          :style="{ backgroundColor: player.color }"
        ></span>
        {{ player.name }}
        <span class="ml-1 text-xs">{{ player.deviceType === 'desktop' ? '💻' : '📱' }}</span>
      </span>
    </span>
  </div>
</template>

<script setup lang="ts">
import type { Player } from '@/shared/types'

defineProps<{
  players: Player[]
}>()
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/host/components/PlayerList.vue
git commit -m "feat: implement host player list"
```

---

## Task 11: Implement Host Main Page

**Files:**
- Create: `frontend/src/host/HostApp.vue`

- [ ] **Step 1: Create HostApp component**

Create `frontend/src/host/HostApp.vue`:
```vue
<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <!-- Header -->
    <div class="flex justify-between items-center px-6 py-4 bg-slate-800">
      <div class="text-sm text-slate-400">{{ serverUrl }}</div>
      <div class="bg-slate-700 px-4 py-2 rounded-lg text-sm">
        📱 二维码（手机扫码加入）
      </div>
    </div>

    <!-- Buzzer Display -->
    <BuzzerDisplay
      :game-state="gameState"
      :winner="winner"
    />

    <!-- Control Panel -->
    <ControlPanel
      :game-state="gameState"
      @start="handleStart"
      @reset="handleReset"
    />

    <!-- Player List -->
    <PlayerList :players="players" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'
import type { Player } from '@/shared/types'
import BuzzerDisplay from './components/BuzzerDisplay.vue'
import ControlPanel from './components/ControlPanel.vue'
import PlayerList from './components/PlayerList.vue'

const gameState = ref<'waiting' | 'ready' | 'locked'>('waiting')
const players = ref<Player[]>([])
const winner = ref<Player | null>(null)

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

const serverUrl = computed(() => {
  return window.location.host
})

let ws: WebSocketClient | null = null

onMounted(() => {
  ws = new WebSocketClient(wsUrl.value)

  ws.on('state_changed', (payload) => {
    gameState.value = payload.state
  })

  ws.on('player_list', (payload) => {
    players.value = payload.players
  })

  ws.on('buzz_result', (payload) => {
    if (payload.winner) {
      winner.value = payload.winner
    } else {
      winner.value = null
    }
  })

  ws.connect()
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
})

const handleStart = () => {
  if (ws) {
    ws.send('start_game', {})
  }
}

const handleReset = () => {
  if (ws) {
    ws.send('reset_game', {})
    winner.value = null
  }
}
</script>
```

- [ ] **Step 2: Update App.vue for routing**

Modify `frontend/src/App.vue`:
```vue
<template>
  <HostApp v-if="isHost" />
  <div v-else class="w-full h-full bg-slate-900 text-white flex items-center justify-center">
    <div class="text-2xl">Player App (Coming soon)</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import HostApp from './host/HostApp.vue'

const isHost = ref(false)

onMounted(() => {
  const userAgent = navigator.userAgent.toLowerCase()
  const isMobile = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(userAgent)
  const path = window.location.pathname

  isHost.value = !isMobile && (path === '/' || path === '/host')
})
</script>
```

- [ ] **Step 3: Test host page**

Run: `cd frontend && npm run dev`

Visit: http://localhost:5173

Expected: See host interface with control buttons

- [ ] **Step 4: Commit**

```bash
git add frontend/src/host/ frontend/src/App.vue
git commit -m "feat: implement host main page with routing"
```

---

## Task 12: Implement Player Name Input Screen

**Files:**
- Create: `frontend/src/player/components/NameInput.vue`

- [ ] **Step 1: Create NameInput component**

Create `frontend/src/player/components/NameInput.vue`:
```vue
<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col items-center justify-center p-8">
    <div class="text-6xl mb-8">🎯</div>
    <div class="text-3xl font-bold mb-4">输入你的名字</div>
    <div class="text-slate-400 mb-8">将在抢答时显示</div>

    <input
      v-model="name"
      @keyup.enter="handleSubmit"
      type="text"
      placeholder="例如：张三"
      class="w-full max-w-md px-6 py-4 text-xl bg-slate-800 border-2 border-slate-700 rounded-lg text-white text-center focus:outline-none focus:border-blue-500"
      maxlength="20"
    />

    <button
      @click="handleSubmit"
      :disabled="!name.trim()"
      class="w-full max-w-md mt-6 px-8 py-4 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-700 disabled:cursor-not-allowed rounded-lg text-xl font-semibold transition-colors"
    >
      进入抢答
    </button>

    <div v-if="error" class="mt-4 text-red-400">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const name = ref('')
const error = ref('')

const handleSubmit = () => {
  const trimmed = name.value.trim()

  if (!trimmed) {
    error.value = '请输入名字'
    return
  }

  if (trimmed.length < 2) {
    error.value = '名字至少需要2个字符'
    return
  }

  if (trimmed.length > 20) {
    error.value = '名字最多20个字符'
    return
  }

  // Store name for next screen
  sessionStorage.setItem('playerName', trimmed)

  // Navigate to buzzer screen
  router.push({ path: '/player', query: { step: 'buzzer' } })
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/player/components/NameInput.vue
git commit -m "feat: implement player name input screen"
```

---

## Task 13: Implement Player Buzzer Button

**Files:**
- Create: `frontend/src/player/components/BuzzerButton.vue`

- [ ] **Step 1: Create BuzzerButton component**

Create `frontend/src/player/components/BuzzerButton.vue`:
```vue
<template>
  <div class="w-full h-full bg-slate-900 text-white flex flex-col">
    <!-- Status Header -->
    <div class="text-center py-6 bg-slate-800">
      <div class="text-2xl font-bold">👋 {{ playerName }}</div>
      <div class="mt-2">
        <span
          class="inline-block px-4 py-1 rounded-full text-sm"
          :class="connected ? 'bg-green-500' : 'bg-red-500'"
        >
          {{ connected ? '✓ 已连接' : '✗ 未连接' }}
        </span>
      </div>
    </div>

    <!-- Status Display -->
    <div class="flex-1 flex items-center justify-center px-8">
      <div class="w-full text-center py-16 bg-slate-800 rounded-2xl">
        <div class="text-4xl font-bold">
          {{ statusText }}
        </div>
      </div>
    </div>

    <!-- Buzzer Button -->
    <div class="p-8">
      <button
        @click="handleBuzz"
        @keydown.space.prevent="handleBuzz"
        @keydown.enter.prevent="handleBuzz"
        :disabled="!canBuzz"
        class="w-full h-32 rounded-2xl font-bold text-3xl transition-all disabled:cursor-not-allowed"
        :class="canBuzz ? 'bg-red-500 hover:bg-red-600 active:scale-95' : 'bg-gray-700 opacity-50'"
      >
        抢答！
      </button>
      <div class="text-center mt-4 text-slate-400 text-sm">
        按空格键 或 点击按钮
      </div>
    </div>

    <!-- Early Buzz Modal -->
    <div
      v-if="showEarlyBuzzModal"
      class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center p-8 z-50"
    >
      <div class="bg-slate-800 rounded-2xl p-8 max-w-md w-full text-center">
        <div class="text-5xl mb-4">⚠️</div>
        <div class="text-2xl font-bold mb-4">请等待主持人开始！</div>
        <div class="text-slate-400 mb-6">抢跑会受到时间惩罚</div>
        <button
          @click="closeEarlyBuzzModal"
          class="w-full px-8 py-4 bg-red-500 hover:bg-red-600 rounded-lg text-xl font-semibold transition-colors"
        >
          我知道了
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { WebSocketClient } from '@/shared/websocket'

const props = defineProps<{
  playerName: string
}>()

const gameState = ref<'waiting' | 'ready' | 'locked'>('waiting')
const connected = ref(false)
const showEarlyBuzzModal = ref(false)

const wsUrl = computed(() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
})

const statusText = computed(() => {
  switch (gameState.value) {
    case 'waiting':
      return '⏳ 等待主持人开始...'
    case 'ready':
      return '🚀 准备抢答！'
    case 'locked':
      return '🔒 本轮已结束'
  }
})

const canBuzz = computed(() => {
  return gameState.value === 'ready' && connected.value
})

let ws: WebSocketClient | null = null
const deviceType = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(navigator.userAgent)
  ? 'mobile'
  : 'desktop'

onMounted(() => {
  ws = new WebSocketClient(wsUrl.value)

  ws.on('state_changed', (payload) => {
    gameState.value = payload.state
  })

  ws.on('early_buzz_warning', () => {
    showEarlyBuzzModal.value = true
  })

  ws.on('open', () => {
    connected.value = true

    // Send join message
    ws?.send('join', {
      name: props.playerName,
      deviceType: deviceType
    })
  })

  ws.connect()

  // Global keyboard listener
  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  if (ws) {
    ws.disconnect()
  }
  window.removeEventListener('keydown', handleGlobalKeydown)
})

const handleGlobalKeydown = (e: KeyboardEvent) => {
  if (e.code === 'Space' || e.code === 'Enter') {
    if (canBuzz.value) {
      handleBuzz()
    }
  }
}

const handleBuzz = () => {
  if (!canBuzz.value || !ws) {
    return
  }

  ws.send('buzz', {
    timestamp: Date.now()
  })
}

const closeEarlyBuzzModal = () => {
  showEarlyBuzzModal.value = false
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/player/components/BuzzerButton.vue
git commit -m "feat: implement player buzzer button"
```

---

## Task 14: Implement Player Main Page

**Files:**
- Create: `frontend/src/player/PlayerApp.vue`

- [ ] **Step 1: Create PlayerApp component**

Create `frontend/src/player/PlayerApp.vue`:
```vue
<template>
  <NameInput v-if="currentStep === 'input'" />
  <BuzzerButton v-else :player-name="playerName" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import NameInput from './components/NameInput.vue'
import BuzzerButton from './components/BuzzerButton.vue'

const route = useRoute()
const currentStep = ref('input')
const playerName = ref('')

onMounted(() => {
  // Check if name is already stored
  const storedName = sessionStorage.getItem('playerName')
  if (storedName) {
    playerName.value = storedName
    currentStep.value = 'buzzer'
  }

  // Check query param
  if (route.query.step === 'buzzer' && storedName) {
    currentStep.value = 'buzzer'
  }
})
</script>
```

- [ ] **Step 2: Update App.vue for player routing**

Modify `frontend/src/App.vue`:
```vue
<template>
  <HostApp v-if="isHost" />
  <PlayerApp v-else />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import HostApp from './host/HostApp.vue'
import PlayerApp from './player/PlayerApp.vue'

const isHost = ref(false)

onMounted(() => {
  const userAgent = navigator.userAgent.toLowerCase()
  const isMobile = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(userAgent)
  const path = window.location.pathname

  isHost.value = !isMobile && (path === '/' || path === '/host')
})
</script>
```

- [ ] **Step 3: Test player page**

Run: `cd frontend && npm run dev`

Visit: http://localhost:5173/player

Expected: See name input screen, then buzzer button

- [ ] **Step 4: Commit**

```bash
git add frontend/src/player/
git commit -m "feat: implement player main page"
```

---

## Task 15: Embed Frontend in Go Backend

**Files:**
- Modify: `backend/main.go`

- [ ] **Step 1: Build frontend**

Run: `cd frontend && npm run build`

Expected: `frontend/dist` directory created with built files

- [ ] **Step 2: Update main.go to serve frontend**

Modify `backend/main.go`:
```go
package main

import (
    "embed"
    "io/fs"
    "log"

    "github.com/yockii/lan-buzzer/backend/game"
    "github.com/yockii/lan-buzzer/backend/websocket"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
)

//go:embed all:../frontend/dist
var frontendFS embed.FS

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })

    app.Use(logger.New())
    app.Use(recover.New())

    gameServer := game.NewGameServer()
    wsHandler := websocket.NewHandler(gameServer)
    wsHandler.RegisterRoutes(app)

    // Serve frontend static files
    frontendDist, err := fs.Sub(frontendFS, "frontend/dist")
    if err != nil {
        log.Fatal(err)
    }

    app.Static("/", "/", fiber.Static{
        FileSystem: http.FS(frontendDist),
        Browse:     false,
    })

    log.Println("Server starting on http://localhost:3000")
    log.Fatal(app.Listen(":3000"))
}
```

- [ ] **Step 3: Add missing imports**

Modify `backend/main.go` imports:
```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"

    "github.com/yockii/lan-buzzer/backend/game"
    "github.com/yockii/lan-buzzer/backend/websocket"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
)
```

- [ ] **Step 4: Fix frontend embed path issue**

The embed directive needs the correct relative path. Since main.go is in `backend/`, we need:

Modify `backend/main.go`:
```go
//go:embed all:../frontend/dist
var frontendFS embed.FS
```

But this won't work with embed. Let's create a different approach:

Create `backend/embed/dist` symlink or copy files, or restructure. Better approach:

Modify `backend/main.go`:
```go
//go:embed frontend/dist
var frontendFS embed.FS
```

And move the dist directory:
```bash
mv frontend/dist backend/embed/frontend
```

Actually, let's use a simpler approach with http.FileServer:

Create `backend/embed/empty.go` to ensure the directory exists:
```go
package embed
```

Then modify the build process to copy the frontend dist to backend/embed/dist.

- [ ] **Step 5: Create build script**

Create `build.sh`:
```bash
#!/bin/bash

echo "Building frontend..."
cd frontend
npm run build
cd ..

echo "Copying frontend to backend..."
mkdir -p backend/embed/dist
rm -rf backend/embed/dist/*
cp -r frontend/dist/* backend/embed/dist/

echo "Building backend..."
cd backend
go build -o ../build/lan-buzzer.exe .
cd ..

echo "Build complete: build/lan-buzzer.exe"
```

- [ ] **Step 6: Update main.go for correct embed path**

Modify `backend/main.go`:
```go
//go:embed embed/dist
var frontendFS embed.FS

func main() {
    // ... existing code ...

    // Serve frontend static files
    frontendDist, err := fs.Sub(frontendFS, "embed/dist")
    if err != nil {
        log.Fatal(err)
    }

    app.Static("/", "", fiber.Static{
        FileSystem: http.FS(frontendDist),
        Browse:     false,
    })

    // ... rest of code ...
}
```

- [ ] **Step 7: Test build**

Run: `chmod +x build.sh && ./build.sh`

Expected: Binary created at `build/lan-buzzer.exe`

- [ ] **Step 8: Test binary**

Run: `./build/lan-buzzer.exe`

Visit: http://localhost:3000

Expected: See host interface

- [ ] **Step 9: Commit**

```bash
git add build.sh backend/embed/ backend/main.go
git commit -m "feat: embed frontend in Go binary"
```

---

## Task 16: Add Auto-Browser Launch

**Files:**
- Modify: `backend/main.go`

- [ ] **Step 1: Add browser launch for Windows**

Modify `backend/main.go`:
```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
    "os/exec"
    "runtime"
    "time"

    "github.com/yockii/lan-buzzer/backend/game"
    "github.com/yockii/lan-buzzer/backend/websocket"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/recover"
)

//go:embed embed/dist
var frontendFS embed.FS

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })

    app.Use(logger.New())
    app.Use(recover.New())

    gameServer := game.NewGameServer()
    wsHandler := websocket.NewHandler(gameServer)
    wsHandler.RegisterRoutes(app)

    // Serve frontend static files
    frontendDist, err := fs.Sub(frontendFS, "embed/dist")
    if err != nil {
        log.Fatal(err)
    }

    app.Static("/", "", fiber.Static{
        FileSystem: http.FS(frontendDist),
        Browse:     false,
    })

    // Start server in background
    go func() {
        log.Println("Server starting on http://localhost:3000")
        if err := app.Listen(":3000"); err != nil {
            log.Fatal(err)
        }
    }()

    // Give server time to start
    time.Sleep(500 * time.Millisecond)

    // Open browser
    openBrowser("http://localhost:3000")

    // Keep program running
    select {}
}

func openBrowser(url string) {
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
    case "darwin":
        cmd = exec.Command("open", url)
    case "linux":
        cmd = exec.Command("xdg-open", url)
    default:
        log.Printf("Unsupported platform for auto-browser launch")
        return
    }

    if err := cmd.Start(); err != nil {
        log.Printf("Failed to open browser: %v", err)
    }
}
```

- [ ] **Step 2: Test browser launch**

Run: `./build/lan-buzzer.exe`

Expected: Browser opens automatically to http://localhost:3000

- [ ] **Step 3: Commit**

```bash
git add backend/main.go
git commit -m "feat: add auto-browser launch on startup"
```

---

## Task 17: Add Start/Reset Game WebSocket Endpoints

**Files:**
- Modify: `backend/websocket/handler.go`

- [ ] **Step 1: Add game control message handling**

Modify `backend/websocket/handler.go` handleConnection method:
```go
func (h *Handler) handleConnection(conn *websocket.Conn) {
    playerID := uuid.New().String()
    var player *game.Player

    defer func() {
        conn.Close()
        h.mutex.Lock()
        delete(h.clients, playerID)
        h.mutex.Unlock()

        if player != nil {
            h.server.RemovePlayer(playerID)
            h.broadcastPlayerList()
        }
        log.Printf("Player %s disconnected", playerID)
    }()

    log.Printf("New connection: %s", playerID)

    for {
        msg := new(Message)
        if err := conn.ReadJSON(msg); err != nil {
            log.Printf("Read error: %v", err)
            break
        }

        switch msg.Type {
        case "join":
            h.handleJoin(conn, playerID, msg)
        case "buzz":
            h.handleBuzz(conn, playerID)
        case "start_game":
            h.handleStartGame(conn)
        case "reset_game":
            h.handleResetGame(conn)
        }
    }
}

func (h *Handler) handleStartGame(conn *websocket.Conn) {
    h.StartGame()
}

func (h *Handler) handleResetGame(conn *websocket.Conn) {
    h.ResetGame()
}
```

- [ ] **Step 2: Test game control**

Run: `./build/lan-buzzer.exe`

Click "开始抢答" and "下一题" buttons

Expected: Game state changes broadcast to all clients

- [ ] **Step 3: Commit**

```bash
git add backend/websocket/handler.go
git commit -m "feat: add start/reset game WebSocket endpoints"
```

---

## Task 18: Add Player Removal Feature

**Files:**
- Modify: `backend/websocket/handler.go`
- Modify: `frontend/src/host/components/PlayerList.vue`

- [ ] **Step 1: Add remove player message handler**

Add to `backend/websocket/handler.go`:
```go
func (h *Handler) handleRemovePlayer(conn *websocket.Conn, msg *Message) {
    payload, ok := msg.Payload.(map[string]any)
    if !ok {
        return
    }

    playerID, _ := payload["playerId"].(string)
    if playerID == "" {
        return
    }

    h.server.RemovePlayer(playerID)
    h.broadcastPlayerList()

    // Disconnect the player
    h.mutex.Lock()
    if playerConn, exists := h.clients[playerID]; exists {
        playerConn.Close()
    }
    h.mutex.Unlock()
}
```

Update handleConnection switch:
```go
switch msg.Type {
case "join":
    h.handleJoin(conn, playerID, msg)
case "buzz":
    h.handleBuzz(conn, playerID)
case "start_game":
    h.handleStartGame(conn)
case "reset_game":
    h.handleResetGame(conn)
case "remove_player":
    h.handleRemovePlayer(conn, msg)
}
```

- [ ] **Step 2: Update PlayerList component**

Modify `frontend/src/host/components/PlayerList.vue`:
```vue
<template>
  <div class="text-center py-4 text-slate-400">
    <span class="text-sm">已连接 {{ players.length }} 位选手</span>
    <div class="mt-2 flex flex-wrap justify-center gap-2">
      <span
        v-for="player in players"
        :key="player.id"
        class="inline-flex items-center bg-slate-800 px-3 py-1 rounded-full"
      >
        <span
          class="w-3 h-3 rounded-full mr-2"
          :style="{ backgroundColor: player.color }"
        ></span>
        {{ player.name }}
        <span class="ml-1 text-xs">{{ player.deviceType === 'desktop' ? '💻' : '📱' }}</span>
        <button
          @click="$emit('remove-player', player.id)"
          class="ml-2 text-slate-500 hover:text-red-400 transition-colors"
        >
          ×
        </button>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Player } from '@/shared/types'

defineProps<{
  players: Player[]
}>()

defineEmits<{
  'remove-player': [playerId: string]
}>()
</script>
```

- [ ] **Step 3: Update HostApp to handle remove**

Modify `frontend/src/host/HostApp.vue`:
```vue
<template>
  <!-- ... existing template ... -->
  <PlayerList
    :players="players"
    @remove-player="handleRemovePlayer"
  />
</template>

<script setup lang="ts">
// ... existing imports ...

const handleRemovePlayer = (playerId: string) => {
  if (ws && confirm('确定要移除这位选手吗？')) {
    ws.send('remove_player', { playerId })
  }
}

// ... rest of code ...
</script>
```

- [ ] **Step 4: Test player removal**

Run: `./build/lan-buzzer.exe`

Connect multiple players, click × to remove

Expected: Player removed from list and disconnected

- [ ] **Step 5: Commit**

```bash
git add backend/websocket/handler.go frontend/src/host/components/PlayerList.vue frontend/src/host/HostApp.vue
git commit -m "feat: add player removal feature"
```

---

## Task 19: Polish and Final Testing

**Files:**
- Create: `README.md`
- Create: `.gitignore`

- [ ] **Step 1: Create README.md**

Create `README.md`:
```markdown
# LAN Buzzer

局域网抢答系统 - 支持 5-10 人实时抢答的简单工具

## 功能特点

- 🎯 实时抢答，毫秒级响应
- 👥 支持 5-10 人同时参与
- 💻 电脑、手机均可使用
- 🎨 颜色标识，支持重名
- ⚡ 单文件运行，无需安装依赖
- 🔒 局域网通信，安全可靠

## 使用方法

1. 下载 `lan-buzzer.exe`（Windows）或 `lan-buzzer`（Linux/Mac）
2. 双击运行
3. 浏览器自动打开主持人界面
4. 选手通过二维码或 URL 加入
5. 点击"开始抢答"开始游戏

## 系统要求

- Windows 10+ / Linux / macOS
- 同一 Wi-Fi 局域网
- 现代浏览器（Chrome、Firefox、Safari）

## 技术栈

- 后端: Go 1.25+ + Fiber v3 + WebSocket
- 前端: Vue 3 + TypeScript + TailwindCSS
```

- [ ] **Step 2: Create .gitignore**

Create `.gitignore`:
```
# Binaries
build/
*.exe
*.dll
*.so
*.dylib
lan-buzzer

# Frontend
frontend/node_modules/
frontend/dist/

# Go
backend/embed/dist/
*.sum

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Superpowers
.superpowers/
```

- [ ] **Step 3: End-to-end test**

1. Start server
2. Connect 3 players (2 desktop, 1 mobile)
3. Test all game states
4. Test buzz timing
5. Test early buzz warning
6. Test player removal
7. Test reconnection

- [ ] **Step 4: Build production binary**

Run: `./build.sh`

Expected: `build/lan-buzzer.exe` created

- [ ] **Step 5: Commit**

```bash
git add README.md .gitignore
git commit -m "docs: add README and .gitignore"
```

---

## Task 20: Final Review and Documentation

**Files:**
- Update: `CLAUDE.md`

- [ ] **Step 1: Create CLAUDE.md**

Create `CLAUDE.md`:
```markdown
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

### Backend (Go + Fiber v3)

- `backend/main.go` - Entry point, serves frontend
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
```

- [ ] **Step 2: Verify all tests pass**

Run: `cd backend && go test ./... -v`

Expected: All tests pass

- [ ] **Step 3: Final commit**

```bash
git add CLAUDE.md
git commit -m "docs: add CLAUDE.md with project guidance"
```

- [ ] **Step 4: Create git tag**

Run: `git tag -a v1.0.0 -m "Initial release of LAN Buzzer"`

---

## Summary

This implementation plan covers:

1. ✅ Go backend with Fiber v3 and WebSocket
2. ✅ Vue 3 frontend with TypeScript and TailwindCSS
3. ✅ Real-time game state management
4. ✅ Player color assignment system
5. ✅ Host control panel with start/reset
6. ✅ Player buzzer interface with early warning
7. ✅ Single binary deployment with embedded frontend
8. ✅ Auto-browser launch on startup
9. ✅ Player removal functionality
10. ✅ Complete documentation

The plan follows TDD, YAGNI, and DRY principles with frequent commits and bite-sized tasks.
