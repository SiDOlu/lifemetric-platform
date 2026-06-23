package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/lifemetric/platform/backend/internal/config"
	"github.com/lifemetric/platform/backend/internal/handlers"
	"github.com/lifemetric/platform/backend/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	// 1. Initialize and run the Global WebSocket broadcast Hub in a background thread
	go handlers.GlobalHub.Run()

	app := fiber.New(fiber.Config{
		AppName: "LifeMetrics Ambient Ingestion API Gateway v1.0",
	})

	// Add request logger middleware
	app.Use(logger.New())

	// Health Check Route (Unauthenticated)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"app":     "LifeMetrics Gateway Engine",
			"version": "1.0.0",
		})
	})

	// 2. Real-Time WebSocket Route Upgrader and Ingestion Endpoint
	app.Use("/ws", handlers.UpgradeWebSocketHandler)

	app.Get("/ws/clinical/alerts", websocket.New(func(c *websocket.Conn) {
		client := &handlers.Client{Conn: c}
		handlers.GlobalHub.Register <- client
		defer func() {
			handlers.GlobalHub.Unregister <- client
		}()

		// Keep connection alive, listening for inbound messages or pings
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))

	// Secure API v1 Router Group (Enforcing JWT verification and role scopes)
	api := app.Group("/api/v1")

	// Device Pairing & Registration Endpoints
	devices := api.Group("/devices", middleware.SecureAuthorize("device:write"))
	devices.Post("/register", handlers.RegisterDeviceHandler)
	devices.Post("/calibrate", handlers.CalibrateDeviceHandler)

	// Stream Telemetry Ingestion Endpoints (requires telemetry:write scope)
	telemetry := api.Group("/events", middleware.SecureAuthorize("telemetry:write"))
	telemetry.Post("/ingest", handlers.IngestTelemetryHandler)

	log.Printf("Initializing LifeMetrics API Gateway on port %s...", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Fatal gateway server crash: %v", err)
	}
}
