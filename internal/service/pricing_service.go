package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/samkhedikar/tou-api/internal/model"
	"github.com/samkhedikar/tou-api/internal/repository"
)

var ErrNoPeriodMatch = errors.New("no pricing period matches the given time")

type PricingService struct {
	chargerRepo repository.ChargerRepository
	pricingRepo repository.PricingRepository
}

func NewPricingService(chargerRepo repository.ChargerRepository, pricingRepo repository.PricingRepository) *PricingService {
	return &PricingService{
		chargerRepo: chargerRepo,
		pricingRepo: pricingRepo,
	}
}

// SetPricing replaces all pricing periods for a charger atomically.
func (s *PricingService) SetPricing(ctx context.Context, chargerID uuid.UUID, req model.SetPricingRequest) (model.PricingResponse, error) {
	charger, err := s.chargerRepo.FindByID(ctx, chargerID)
	if err != nil {
		return model.PricingResponse{}, err
	}

	if err := validatePeriods(req.Periods); err != nil {
		return model.PricingResponse{}, err
	}

	// Build period models.
	periods := make([]model.PricingPeriod, 0, len(req.Periods))
	for _, p := range req.Periods {
		periods = append(periods, model.PricingPeriod{
			ID:          uuid.New(),
			ChargerID:   chargerID,
			StartHour:   p.StartHour,
			StartMinute: p.StartMinute,
			EndHour:     p.EndHour,
			EndMinute:   p.EndMinute,
			PricePerKWh: p.PricePerKWh,
			CreatedAt:   time.Now().UTC(),
		})
	}

	if err := s.pricingRepo.ReplaceAll(ctx, chargerID, periods); err != nil {
		return model.PricingResponse{}, err
	}

	return model.PricingResponse{
		ChargerID:   charger.ID,
		ChargerName: charger.Name,
		Timezone:    charger.Timezone,
		Periods:     periods,
	}, nil
}

// GetPricing returns all pricing periods for a charger.
func (s *PricingService) GetPricing(ctx context.Context, chargerID uuid.UUID) (model.PricingResponse, error) {
	charger, err := s.chargerRepo.FindByID(ctx, chargerID)
	if err != nil {
		return model.PricingResponse{}, err
	}

	periods, err := s.pricingRepo.FindByChargerID(ctx, chargerID)
	if err != nil {
		return model.PricingResponse{}, err
	}

	return model.PricingResponse{
		ChargerID:   charger.ID,
		ChargerName: charger.Name,
		Timezone:    charger.Timezone,
		Periods:     periods,
	}, nil
}

// GetPriceAt returns the applicable TOU price for a charger at a given UTC time.
// Converts to the charger's local timezone before matching against periods.
func (s *PricingService) GetPriceAt(ctx context.Context, chargerID uuid.UUID, at time.Time) (model.GetPriceAtResponse, error) {
	charger, err := s.chargerRepo.FindByID(ctx, chargerID)
	if err != nil {
		return model.GetPriceAtResponse{}, err
	}

	loc, err := time.LoadLocation(charger.Timezone)
	if err != nil {
		return model.GetPriceAtResponse{}, fmt.Errorf("load timezone %q: %w", charger.Timezone, err)
	}

	localTime := at.In(loc)
	localMinutes := localTime.Hour()*60 + localTime.Minute()

	log.Printf("GetPriceAt: UTC=%s → local=%s (%s) localMinutes=%d",
		at.UTC().Format(time.RFC3339),
		localTime.Format("2006-01-02 15:04:05"),
		charger.Timezone,
		localMinutes,
	)

	periods, err := s.pricingRepo.FindByChargerID(ctx, chargerID)
	if err != nil {
		return model.GetPriceAtResponse{}, err
	}

	if len(periods) == 0 {
		return model.GetPriceAtResponse{}, ErrNoPeriodMatch
	}

	for _, p := range periods {
		startMin := p.StartHour*60 + p.StartMinute
		endMin := p.EndHour*60 + p.EndMinute
		if localMinutes >= startMin && localMinutes < endMin {
			return model.GetPriceAtResponse{
				ChargerID:   charger.ID,
				ChargerName: charger.Name,
				QueryTime:   at.UTC().Format(time.RFC3339),
				LocalTime:   localTime.Format("15:04"),
				PricePerKWh: p.PricePerKWh,
				PeriodStart: fmt.Sprintf("%02d:%02d", p.StartHour, p.StartMinute),
				PeriodEnd:   fmt.Sprintf("%02d:%02d", p.EndHour, p.EndMinute),
				Timezone:    charger.Timezone,
			}, nil
		}
	}

	return model.GetPriceAtResponse{}, ErrNoPeriodMatch
}

// BulkSetPricing applies the same pricing to multiple chargers.
func (s *PricingService) BulkSetPricing(ctx context.Context, req model.BulkSetPricingRequest) (model.BulkSetPricingResponse, error) {
	if err := validatePeriods(req.Periods); err != nil {
		return model.BulkSetPricingResponse{}, err
	}

	resp := model.BulkSetPricingResponse{Errors: make(map[string]string)}
	for _, cid := range req.ChargerIDs {
		_, err := s.SetPricing(ctx, cid, model.SetPricingRequest{Periods: req.Periods})
		if err != nil {
			resp.Errors[cid.String()] = err.Error()
			continue
		}
		resp.Updated = append(resp.Updated, cid)
	}
	return resp, nil
}

// validatePeriods checks end > start and no overlapping periods.
func validatePeriods(periods []model.PricingPeriodInput) error {
	for i, p := range periods {
		if p.EndHour*60+p.EndMinute <= p.StartHour*60+p.StartMinute {
			return fmt.Errorf("period %d: end time must be after start time", i)
		}
	}
	for i := 0; i < len(periods); i++ {
		for j := i + 1; j < len(periods); j++ {
			aS := periods[i].StartHour*60 + periods[i].StartMinute
			aE := periods[i].EndHour*60 + periods[i].EndMinute
			bS := periods[j].StartHour*60 + periods[j].StartMinute
			bE := periods[j].EndHour*60 + periods[j].EndMinute
			if aS < bE && bS < aE {
				return fmt.Errorf("periods %d and %d overlap", i, j)
			}
		}
	}
	return nil
}
