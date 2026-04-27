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
