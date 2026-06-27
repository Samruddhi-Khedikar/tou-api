package model

import (
	"time"

	"github.com/google/uuid"
)

type Charger struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PricingPeriod struct {
	ID          uuid.UUID `json:"id"`
	ChargerID   uuid.UUID `json:"charger_id"`
	StartHour   int       `json:"start_hour"`
	StartMinute int       `json:"start_minute"`
	EndHour     int       `json:"end_hour"`
	EndMinute   int       `json:"end_minute"`
	PricePerKWh float64   `json:"price_per_kwh"`
	CreatedAt   time.Time `json:"created_at"`
}

// --- Request types ---

type CreateChargerRequest struct {
	Name     string `json:"name"     validate:"required"`
	Location string `json:"location" validate:"required"`
	Timezone string `json:"timezone" validate:"required"`
}

type PricingPeriodInput struct {
	StartHour   int     `json:"start_hour"    validate:"min=0,max=23"`
	StartMinute int     `json:"start_minute"  validate:"min=0,max=59"`
	EndHour     int     `json:"end_hour"      validate:"min=0,max=24"`
	EndMinute   int     `json:"end_minute"    validate:"min=0,max=59"`
	PricePerKWh float64 `json:"price_per_kwh" validate:"gte=0"`
}

type SetPricingRequest struct {
	Periods []PricingPeriodInput `json:"periods" validate:"required,min=1,dive"`
}

type BulkSetPricingRequest struct {
	ChargerIDs []uuid.UUID          `json:"charger_ids" validate:"required,min=1"`
	Periods    []PricingPeriodInput `json:"periods"     validate:"required,min=1,dive"`
}

// --- Response types ---

type PricingResponse struct {
	ChargerID   uuid.UUID       `json:"charger_id"`
	ChargerName string          `json:"charger_name"`
	Timezone    string          `json:"timezone"`
	Periods     []PricingPeriod `json:"periods"`
}

type BulkSetPricingResponse struct {
	Updated []uuid.UUID       `json:"updated"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type GetPriceAtResponse struct {
	ChargerID   uuid.UUID `json:"charger_id"`
	ChargerName string    `json:"charger_name"`
	QueryTime   string    `json:"query_time"`
	LocalTime   string    `json:"local_time"`
	PricePerKWh float64   `json:"price_per_kwh"`
	PeriodStart string    `json:"period_start"`
	PeriodEnd   string    `json:"period_end"`
	Timezone    string    `json:"timezone"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
