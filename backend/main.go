package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/yockii/lan_qr/backend/game"
	"github.com/yockii/lan_qr/backend/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
