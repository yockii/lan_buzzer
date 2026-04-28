package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/yockii/lan_qr/backend/game"
	"github.com/yockii/lan_qr/backend/quiz"
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

	// Try to load question bank
	questionBankPath := "questions.txt"
	if qb, err := quiz.LoadQuestionBank(questionBankPath); err == nil {
		// Convert quiz.QuestionBank to game.QuestionBank interface
		gameServer.SetQuestionBank(qb)
		log.Printf("Loaded question bank, mode: quiz")
	} else {
		log.Printf("No question bank found (%v), mode: buzzer", err)
	}

	wsHandler := websocket.NewHandler(gameServer)
	wsHandler.RegisterRoutes(app)

	// API endpoint to get server info
	app.Get("/api/info", func(c *fiber.Ctx) error {
		localIP := getLocalIP()
		allIPs := getAllLocalIPs()
		mode := "buzzer"
		if gameServer.HasQuestionBank() {
			mode = "quiz"
		}
		return c.JSON(fiber.Map{
			"serverUrl": fmt.Sprintf("http://%s:3000", localIP),
			"localIP":  localIP,
			"allIPs":   allIPs,
			"mode":     mode,
		})
	})

	// Serve frontend static files
	frontendDist, err := fs.Sub(frontendFS, "embed/dist")
	if err != nil {
		log.Fatal(err)
	}

	// Handle static file serving
	app.Get("/*", func(c *fiber.Ctx) error {
		requestPath := c.Params("*")

		// Default to index.html for root or empty path
		if requestPath == "" || requestPath == "/" {
			requestPath = "index.html"
		}

		// Try to read the requested file
		data, err := fs.ReadFile(frontendDist, strings.TrimPrefix(requestPath, "/"))
		if err != nil {
			// Fall back to index.html for SPA routing
			data, err = fs.ReadFile(frontendDist, "index.html")
			if err != nil {
				return c.Status(404).SendString("404 - Page Not Found")
			}
		}

		// Set proper Content-Type header
		c.Set("Content-Type", getContentType(requestPath))

		// Return the file content
		return c.Send(data)
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

func getLocalIP() string {
	// Get all network interfaces
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Failed to get interfaces: %v", err)
		return "localhost"
	}

	// Log all addresses for debugging
	log.Println("All network addresses:")
	for _, addr := range addrs {
		log.Printf("  - %v", addr)
	}

	// First pass: look for 192.168.x.x addresses (LAN)
	for _, addr := range addrs {
		var ipnet *net.IPNet
		var ok bool
		if ipnet, ok = addr.(*net.IPNet); ok {
			ip := ipnet.IP.To4()
			// Check if it's not loopback and is 192.168.x.x
			if ip != nil && !ipnet.IP.IsLoopback() && ip[0] == 192 && ip[1] == 168 {
				log.Printf("Selected 192.168.x.x address: %s", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
		_ = ok
	}

	// Second pass: any non-loopback IPv4 address
	for _, addr := range addrs {
		var ipnet *net.IPNet
		var ok bool
		if ipnet, ok = addr.(*net.IPNet); ok {
			if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				log.Printf("Selected fallback address: %s", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
		_ = ok
	}

	// Fallback to localhost
	log.Println("No suitable address found, using localhost")
	return "localhost"
}

func getAllLocalIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Failed to get interfaces: %v", err)
		return []string{"localhost"}
	}

	var ips []string
	for _, addr := range addrs {
		var ipnet *net.IPNet
		var ok bool
		if ipnet, ok = addr.(*net.IPNet); ok {
			ip := ipnet.IP.To4()
			// Filter out loopback, link-local (169.254.x.x), and docker/VM IPs
			if ip != nil && !ipnet.IP.IsLoopback() && ip[0] != 169 && ip[0] != 172 {
				ips = append(ips, ipnet.IP.String())
			}
		}
		_ = ok
	}

	// Always add localhost as fallback
	if len(ips) == 0 {
		ips = append(ips, "localhost")
	}

	return ips
}

func getContentType(path string) string {
	// Extract file extension
	if idx := strings.LastIndex(path, "."); idx != -1 {
		ext := strings.ToLower(path[idx:])
		switch ext {
		case ".html", ".htm":
			return "text/html; charset=utf-8"
		case ".css":
			return "text/css; charset=utf-8"
		case ".js":
			return "application/javascript; charset=utf-8"
		case ".json", ".map":
			return "application/json; charset=utf-8"
		case ".png":
			return "image/png"
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".svg":
			return "image/svg+xml"
		case ".ico":
			return "image/x-icon"
		case ".woff":
			return "font/woff"
		case ".woff2":
			return "font/woff2"
		case ".ttf":
			return "font/ttf"
		default:
			return "text/plain; charset=utf-8"
		}
	}
	return "text/html; charset=utf-8"  // Default to HTML for SPA
}
