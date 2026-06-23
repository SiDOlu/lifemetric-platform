package models

import (
	"time"
)

// Tenant represents a clinical facility or home care network (B2B Tenant)
type Tenant struct {
	ID           string    `json:"id" db:"id"`
	BusinessName string    `json:"business_name" db:"business_name"`
	PlanType     string    `json:"plan_type" db:"plan_type"` // e.g. "free", "premium"
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Device represents a physical wall/ceiling-mounted LifeMetric sensor array
type Device struct {
	UUID            string    `json:"device_uuid" db:"device_uuid"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	RoomLabel       string    `json:"room_label" db:"room_label"`
	DevicePlacement string    `json:"device_placement" db:"device_placement"` // e.g. "ceiling", "wall"
	PublicKey       string    `json:"public_key" db:"public_key"`
	Status          string    `json:"status" db:"status"` // e.g. "registered", "calibrating", "active", "offline"
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// IngestionPayload represents the anonymized behavioral metadata received from Edge-AI sensors
type IngestionPayload struct {
	DeviceUUID string `json:"device_uuid"`
	Timestamp  string `json:"timestamp"`
	Data       struct {
		CenterOfMassX float64 `json:"com_x"`
		CenterOfMassY float64 `json:"com_y"`
		CenterOfMassZ float64 `json:"com_z"`
		VelocityZ     float64 `json:"velocity_z"`
		IsDangling    bool    `json:"is_dangling"`
		IsPresent     bool    `json:"is_present"`
		VitalsBPM     int     `json:"vitals_bpm"`     // Contactless heart rate
		VitalsBreaths int     `json:"vitals_breaths"` // Contactless respiratory rate
	} `json:"data"`
}

// WebhookSubscription represents a partner integration subscription endpoint
type WebhookSubscription struct {
	ID         string    `json:"id" db:"id"`
	TargetURL  string    `json:"target_url" db:"target_url"`
	SecretKey  string    `json:"secret_key" db:"secret_key"`
	EventTypes []string  `json:"event_types" db:"event_types"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// AlertEvent represents critical alerts like falls or bed exits dispatched to caregivers
type AlertEvent struct {
	EventID    string      `json:"event_id"`
	EventType  string      `json:"event_type"` // e.g. "alert.fall.predictive", "alert.fall.confirmed", "alert.bed_exit.dangling"
	Timestamp  time.Time   `json:"timestamp"`
	DeviceUUID string      `json:"device_uuid"`
	RoomLabel  string      `json:"room_label"`
	Data       interface{} `json:"data"`
}
