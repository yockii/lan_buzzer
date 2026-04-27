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
