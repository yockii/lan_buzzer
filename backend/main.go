package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

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

	log.Println("Server starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
