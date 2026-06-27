package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/samkhedikar/tou-api/internal/model"
	"github.com/samkhedikar/tou-api/internal/repository"
	"github.com/samkhedikar/tou-api/internal/service"
)

// PUT /chargers/{chargerID}/pricing
// Replaces all TOU pricing periods for a charger atomically.
// Required: periods (min 1), each with start_hour, end_hour, price_per_kwh
// Validations: no overlapping periods, end time must be after start time
func (h *Handler) SetPricing(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "chargerID"))
	if !ok {
		return
	}

	var req model.SetPricingRequest
	if !h.parseAndValidate(r, w, &req) {
		return
	}

	resp, err := h.pricingSvc.SetPricing(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "charger not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GET /chargers/{chargerID}/pricing
// Returns all pricing periods for a charger.
func (h *Handler) GetPricing(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "chargerID"))
	if !ok {
		return
	}

	resp, err := h.pricingSvc.GetPricing(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "charger not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if len(resp.Periods) == 0 {
		writeError(w, http.StatusNotFound, "no pricing configured for this charger")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GET /chargers/{chargerID}/price?at=2026-06-27T14:30:00Z
// Returns the applicable TOU price for a charger at a given point in time.
// Query param `at` is optional — defaults to current time if not provided.
// `at` must be RFC3339 UTC format e.g. 2026-06-27T14:30:00Z
func (h *Handler) GetPriceAt(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "chargerID"))
	if !ok {
		return
	}

	at := time.Now().UTC()
	if raw := r.URL.Query().Get("at"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid 'at' param; use RFC3339 e.g. 2026-06-27T14:30:00Z")
			return
		}
		at = parsed
	}

	resp, err := h.pricingSvc.GetPriceAt(r.Context(), id, at)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "charger not found")
			return
		}
		if errors.Is(err, service.ErrNoPeriodMatch) {
			writeError(w, http.StatusNotFound, "no pricing period covers the requested time")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// POST /chargers/bulk/pricing
// Applies the same TOU pricing to multiple chargers at once.
// Required: charger_ids (min 1), periods (min 1)
// Partial success supported — returns which chargers succeeded and which failed.
func (h *Handler) BulkSetPricing(w http.ResponseWriter, r *http.Request) {
	var req model.BulkSetPricingRequest
	if !h.parseAndValidate(r, w, &req) {
		return
	}

	resp, err := h.pricingSvc.BulkSetPricing(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	status := http.StatusOK
	if len(resp.Updated) == 0 {
		status = http.StatusMultiStatus
	}
	writeJSON(w, status, resp)
}
