package game

import (
	"testing"
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
