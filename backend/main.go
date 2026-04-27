package main

import (
	"log"

	"github.com/yockii/lan_qr/backend/game"
	"github.com/yockii/lan_qr/backend/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
