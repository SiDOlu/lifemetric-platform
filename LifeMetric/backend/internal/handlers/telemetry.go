package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lifemetric/platform/backend/internal/models"
)

// IngestTelemetryHandler handles async stream ingestion of Edge-AI metadata payloads
func IngestTelemetryHandler(c *fiber.Ctx) error {
	payload := new(models.IngestionPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse telemetry payload",
		})
	}

	// Retrieve corresponding device
	device, exists := ActiveDevices[payload.DeviceUUID]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Unknown device UUID",
		})
	}

	// Validate Tenant context (security boundary)
	authTenantID := c.Locals("tenant_id").(string)
	if device.TenantID != authTenantID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: Device belongs to another facility tenant",
		})
	}

	// 1. Process Kinematics: Fall Detection / Predication
	var alert *models.AlertEvent

	// Check for predicted / incipient fall (Unstable Descent: high downward velocity before impact)
	if payload.Data.VelocityZ < -1.8 {
		alert = &models.AlertEvent{
			EventID:    uuid.New().String(),
			EventType:  "alert.fall.predictive",
			Timestamp:  time.Now(),
			DeviceUUID: device.UUID,
			RoomLabel:  device.RoomLabel,
			Data: map[string]interface{}{
				"predicted_impact_seconds": 1.2,
				"velocity_z_m_s":          payload.Data.VelocityZ,
				"confidence_score":         0.92,
				"calibration_phase":       payload.Data.CalibrationPhase,
				"confidence_status":       payload.Data.ConfidenceStatus,
			},
		}
	} else if payload.Data.IsDangling {
		// Check for dangling state (Intent to exit bed unassisted)
		alert = &models.AlertEvent{
			EventID:    uuid.New().String(),
			EventType:  "alert.bed_exit.dangling",
			Timestamp:  time.Now(),
			DeviceUUID: device.UUID,
			RoomLabel:  device.RoomLabel,
			Data: map[string]interface{}{
				"action":                   "dangling_detected",
				"elapsed_dangling_seconds": 15,
				"calibration_phase":       payload.Data.CalibrationPhase,
				"confidence_status":       payload.Data.ConfidenceStatus,
			},
		}
	}

	// Dispatch alerts asynchronously if triggered
	if alert != nil {
		log.Printf("[ALERT DISPATCHER] High-Priority Event Triggered: %s in %s", alert.EventType, alert.RoomLabel)
		
		// Serialize and broadcast alert to connected web terminals in real-time
		if alertJSON, err := json.Marshal(alert); err == nil {
			GlobalHub.Broadcast <- alertJSON
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"status":             "accepted",
			"message":            "Telemetry received, critical alert dispatched",
			"triggered_event_id": alert.EventID,
			"event_type":         alert.EventType,
		})
	}

	// Under standard conditions, return a fast 202 Accepted (Async ACK)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status":  "accepted",
		"message": "Telemetry metadata ingested successfully",
	})
}
