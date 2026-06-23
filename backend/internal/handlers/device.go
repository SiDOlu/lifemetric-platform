package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lifemetric/platform/backend/internal/models"
)

// In-memory mock database for scaffolding/prototyping
var ActiveDevices = make(map[string]*models.Device)

// RegisterDeviceHandler handles registering physical sensors and pairing them to facilities
func RegisterDeviceHandler(c *fiber.Ctx) error {
	device := new(models.Device)
	if err := c.BodyParser(device); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse device registration payload",
		})
	}

	if device.UUID == "" {
		device.UUID = uuid.New().String()
	}

	// Extract Tenant ID from authenticated token context (guarantees multi-tenant security boundary)
	device.TenantID = c.Locals("tenant_id").(string)
	device.Status = "registered"
	device.CreatedAt = time.Now()

	ActiveDevices[device.UUID] = device

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "registered",
		"device_uuid": device.UUID,
		"room_label":  device.RoomLabel,
		"tenant_id":   device.TenantID,
		"paired_at":   device.CreatedAt,
	})
}

// CalibrateDeviceHandler triggers the 48-hour local room calibration phase
func CalibrateDeviceHandler(c *fiber.Ctx) error {
	payload := struct {
		DeviceUUID string `json:"device_uuid"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse calibration payload",
		})
	}

	device, exists := ActiveDevices[payload.DeviceUUID]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Target device not found",
		})
	}

	// Verify Tenant matches the authenticated user's tenant context (Row-Level Security boundary)
	authTenantID := c.Locals("tenant_id").(string)
	if device.TenantID != authTenantID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: Device belongs to another facility tenant",
		})
	}

	device.Status = "calibrating"

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           "calibrating",
		"device_uuid":      device.UUID,
		"phase_duration":   "48 hours",
		"started_at":       time.Now(),
		"message":          "Edge-AI calibration triggered. Mapping static background clutter.",
	})
}
