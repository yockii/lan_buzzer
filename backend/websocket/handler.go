package websocket

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/yockii/lan_qr/backend/game"
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
	app.Get("/ws", websocket.New(h.handleConnection))
}

func (h *Handler) handleConnection(conn *websocket.Conn) {
	playerID := uuid.New().String()
	player := make(chan *game.Player, 1)

	defer func() {
		conn.Close()
		h.mutex.Lock()
		delete(h.clients, playerID)
		h.mutex.Unlock()

		select {
		case p := <-player:
			h.server.RemovePlayer(p.ID)
			h.broadcastPlayerList()
		default:
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
			h.handleJoin(conn, playerID, msg, player)
		case "buzz":
			h.handleBuzz(conn, playerID)
		case "start_game":
			h.handleStartGame(conn)
		case "reset_game":
			h.handleResetGame(conn)
		case "remove_player":
			h.handleRemovePlayer(conn, msg)
		}
	}
}

func (h *Handler) handleJoin(conn *websocket.Conn, playerID string, msg *Message, playerChan chan<- *game.Player) {
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

	newPlayer := &game.Player{
		ID:         playerID,
		Name:       name,
		Color:      game.GetPlayerColor(playerID, existingColors),
		DeviceType: deviceType,
	}

	h.server.AddPlayer(newPlayer)
	h.broadcastPlayerList()
	h.broadcastState(conn)

	playerChan <- newPlayer
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

func (h *Handler) handleStartGame(conn *websocket.Conn) {
	h.StartGame()
}

func (h *Handler) handleResetGame(conn *websocket.Conn) {
	h.ResetGame()
}

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
