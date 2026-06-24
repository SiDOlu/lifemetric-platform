package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lifemetric/platform/backend/internal/models"
)

func TestIngestTelemetryHandler_WithNewFields(t *testing.T) {
	// Setup
	app := fiber.New(fiber.Config{
		DisableKeepalive: true,
	})
	
	// Start the hub in the background to prevent blocking on Broadcast
	go GlobalHub.Run()

	// Mock device
	deviceUUID := "test-device-uuid"
	tenantID := "test-tenant-id"
	ActiveDevices[deviceUUID] = &models.Device{
		UUID:      deviceUUID,
		TenantID:  tenantID,
		RoomLabel: "Test Room",
	}

	// Mock middleware to set tenant_id
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", tenantID)
		return c.Next()
	})

	// Register handler
	app.Post("/ingest", IngestTelemetryHandler)

	// Prepare payload with new fields
	payload := map[string]interface{}{
		"device_uuid": deviceUUID,
		"timestamp":   "2023-10-27T10:00:00Z",
		"data": map[string]interface{}{
			"com_x":             1.0,
			"com_y":             1.0,
			"com_z":             1.0,
			"velocity_z":        -2.0, // Trigger fall alert
			"is_dangling":       false,
			"is_present":        true,
			"vitals_bpm":        70,
			"vitals_breaths":    18,
			"calibration_phase": "mapping_static_clutter",
			"confidence_status": "high_confidence",
		},
	}
	body, _ := json.Marshal(payload)

	// Execute request
	req := httptest.NewRequest("POST", "/ingest", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)

	// Assertions
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected status 202 Accepted, got %d", resp.StatusCode)
	}
}

func TestIngestTelemetryHandler_DanglingWithNewFields(t *testing.T) {
	// Setup
	app := fiber.New(fiber.Config{
		DisableKeepalive: true,
	})
	
	// Start the hub in the background to prevent blocking on Broadcast
	go GlobalHub.Run()

	// Mock device
	deviceUUID := "test-device-uuid"
	tenantID := "test-tenant-id"
	ActiveDevices[deviceUUID] = &models.Device{
		UUID:      deviceUUID,
		TenantID:  tenantID,
		RoomLabel: "Test Room",
	}

	// Mock middleware to set tenant_id
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("tenant_id", tenantID)
		return c.Next()
	})

	// Register handler
	app.Post("/ingest", IngestTelemetryHandler)

	// Prepare payload with new fields
	payload := map[string]interface{}{
		"device_uuid": deviceUUID,
		"timestamp":   "2023-10-27T10:00:00Z",
		"data": map[string]interface{}{
			"com_x":             1.0,
			"com_y":             1.0,
			"com_z":             1.0,
			"velocity_z":        0.0,
			"is_dangling":       true, // Trigger dangling alert
			"is_present":        true,
			"vitals_bpm":        70,
			"vitals_breaths":    18,
			"calibration_phase": "initial_setup",
			"confidence_status": "low_confidence",
		},
	}
	body, _ := json.Marshal(payload)

	// Execute request
	req := httptest.NewRequest("POST", "/ingest", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)

	// Assertions
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected status 202 Accepted, got %d", resp.StatusCode)
	}
}
