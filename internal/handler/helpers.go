package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samkhedikar/tou-api/internal/model"
)

// ChargerService defines charger business operations.
type ChargerService interface {
	CreateCharger(ctx context.Context, req model.CreateChargerRequest) (model.Charger, error)
	GetCharger(ctx context.Context, id uuid.UUID) (model.Charger, error)
}

// PricingService defines pricing business operations.
type PricingService interface {
	SetPricing(ctx context.Context, chargerID uuid.UUID, req model.SetPricingRequest) (model.PricingResponse, error)
	GetPricing(ctx context.Context, chargerID uuid.UUID) (model.PricingResponse, error)
	GetPriceAt(ctx context.Context, chargerID uuid.UUID, at time.Time) (model.GetPriceAtResponse, error)
	BulkSetPricing(ctx context.Context, req model.BulkSetPricingRequest) (model.BulkSetPricingResponse, error)
}

// Handler holds all service dependencies.
type Handler struct {
	chargerSvc ChargerService
	pricingSvc PricingService
	validate   *validator.Validate
}

func New(chargerSvc ChargerService, pricingSvc PricingService) *Handler {
	return &Handler{
		chargerSvc: chargerSvc,
		pricingSvc: pricingSvc,
		validate:   validator.New(),
	}
}

func (h *Handler) parseAndValidate(r *http.Request, w http.ResponseWriter, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	if err := h.validate.Struct(dst); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func parseUUID(w http.ResponseWriter, raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(raw)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid charger ID")
		return uuid.UUID{}, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, model.ErrorResponse{Error: msg})
}
