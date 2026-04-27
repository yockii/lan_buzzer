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
