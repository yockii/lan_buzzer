package websocket

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/yockii/lan_qr/backend/game"
	"github.com/yockii/lan_qr/backend/quiz"
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
		case "quiz_start":
			h.handleQuizStart(conn)
		case "quiz_answer":
			h.handleQuizAnswer(conn, playerID, msg)
		case "quiz_judge":
			h.handleQuizJudge(conn, msg)
		case "quiz_next":
			h.handleQuizNext(conn)
		case "quiz_reset":
			h.handleQuizReset(conn)
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

	log.Printf("Player join: name=%s, id=%s, device=%s", name, playerID, deviceType)

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

	// 如果是主持人，不加入玩家列表
	if name == "__HOST__" {
		log.Printf("Host connected: %s", playerID)
		h.broadcastState(conn)
		return
	}

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
	log.Printf("Buzz received from %s, state=%s", playerID, state)

	if state == game.StateWaiting {
		log.Printf("Early buzz warning sent to %s", playerID)
		h.sendEarlyBuzzWarning(conn)
		return
	}

	if state == game.StateLocked {
		log.Printf("Game locked, ignoring buzz from %s", playerID)
		return
	}

	if h.server.RecordBuzz(playerID) {
		winner := h.server.GetWinner()
		log.Printf("Buzz result: winner=%s", winner.Name)
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

	log.Printf("Broadcasting message type=%s to %d clients", msg.Type, len(h.clients))
	for _, conn := range h.clients {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error writing to client: %v", err)
		}
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

func (h *Handler) handleQuizStart(conn *websocket.Conn) {
	question := h.server.StartQuizQuestion()
	if question == nil {
		h.broadcastNoQuestions()
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
		log.Printf("[DEBUG] handleQuizJudge: payload is not map[string]any")
		return
	}

	playerID, _ := payload["playerId"].(string)
	if playerID == "" {
		log.Printf("[DEBUG] handleQuizJudge: playerId is empty")
		return
	}

	// Handle both boolean and string types for 'correct'
	var correct bool
	switch v := payload["correct"].(type) {
	case bool:
		correct = v
		log.Printf("[DEBUG] handleQuizJudge: correct is bool: %v", v)
	case string:
		correct = v == "true" || v == "1"
		log.Printf("[DEBUG] handleQuizJudge: correct is string: %s -> %v", v, correct)
	default:
		log.Printf("[DEBUG] handleQuizJudge: correct is unknown type: %T", v)
		return
	}

	log.Printf("[DEBUG] handleQuizJudge: playerId=%s correct=%v", playerID, correct)
	h.server.JudgeAnswer(playerID, correct)

	// Always broadcast the updated answer
	quizState := h.server.GetQuizState()
	if quizState != nil {
		if playerAnswer, exists := quizState.Answers[playerID]; exists {
			h.broadcastQuizAnswerUpdate(playerAnswer)
		}
	}
}

func (h *Handler) handleQuizNext(conn *websocket.Conn) {
	log.Printf("[DEBUG] handleQuizNext: getting next question")
	question := h.server.NextQuizQuestion()
	if question != nil {
		log.Printf("[DEBUG] handleQuizNext: broadcasting question ID=%s", question.ID)
		h.broadcastQuizQuestion(question)
	} else {
		log.Printf("[DEBUG] handleQuizNext: no more questions")
		h.broadcastNoQuestions()
	}
}

func (h *Handler) handleQuizReset(conn *websocket.Conn) {
	if h.server.ResetQuizQuestionBank() {
		log.Printf("Quiz question bank reset")
		// Broadcast reset notification to clear UI
		h.broadcast(Message{
			Type:    "quiz_reset",
			Payload: map[string]any{},
		})
	}
}

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

	log.Printf("[DEBUG] broadcastQuizQuestion: broadcasting quiz_question to %d clients", len(h.clients))
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
		"playerId":    playerAnswer.PlayerID,
		"playerName":  playerName,
		"playerColor":  playerColor,
		"answer":      playerAnswer.Answer,
		"status":      string(playerAnswer.Status),
		"timestamp":   playerAnswer.Timestamp,
	}

	msg := Message{
		Type:    "quiz_answer_update",
		Payload: payload,
	}

	h.broadcast(msg)
}

func (h *Handler) broadcastNoQuestions() {
	msg := Message{
		Type:    "quiz_no_questions",
		Payload: map[string]any{},
	}

	h.broadcast(msg)
}
