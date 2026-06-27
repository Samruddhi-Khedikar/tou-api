package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// Charger routes
	r.Post("/chargers", h.CreateCharger)
	r.Get("/chargers/{chargerID}", h.GetCharger)

	// Pricing routes
	r.Put("/chargers/{chargerID}/pricing", h.SetPricing)
	r.Get("/chargers/{chargerID}/pricing", h.GetPricing)
	r.Get("/chargers/{chargerID}/price", h.GetPriceAt)

	// Bulk pricing
	r.Post("/chargers/bulk/pricing", h.BulkSetPricing)

	return r
}
