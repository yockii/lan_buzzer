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
