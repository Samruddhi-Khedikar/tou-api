package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samkhedikar/tou-api/internal/model"
	"github.com/samkhedikar/tou-api/internal/repository"
)

// POST /chargers
// Creates a new EV charging station.
// Required: name, location, timezone (valid IANA)
func (h *Handler) CreateCharger(w http.ResponseWriter, r *http.Request) {
	var req model.CreateChargerRequest
	if !h.parseAndValidate(r, w, &req) {
		return
	}

	charger, err := h.chargerSvc.CreateCharger(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, charger)
}

// GET /chargers/{chargerID}
// Returns details of a specific charging station.
func (h *Handler) GetCharger(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, chi.URLParam(r, "chargerID"))
	if !ok {
		return
	}

	charger, err := h.chargerSvc.GetCharger(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "charger not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, charger)
}
