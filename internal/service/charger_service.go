package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samkhedikar/tou-api/internal/model"
	"github.com/samkhedikar/tou-api/internal/repository"
)

type ChargerService struct {
	repo repository.ChargerRepository
}

func NewChargerService(repo repository.ChargerRepository) *ChargerService {
	return &ChargerService{repo: repo}
}

func (s *ChargerService) CreateCharger(ctx context.Context, req model.CreateChargerRequest) (model.Charger, error) {
	// Validate IANA timezone — can't be done via struct tag.
	if _, err := time.LoadLocation(req.Timezone); err != nil {
		return model.Charger{}, fmt.Errorf("invalid timezone %q: %w", req.Timezone, err)
	}

	charger := model.Charger{
		ID:       uuid.New(),
		Name:     req.Name,
		Location: req.Location,
		Timezone: req.Timezone,
	}

	if err := s.repo.Insert(ctx, charger); err != nil {
		return model.Charger{}, err
	}

	return s.repo.FindByID(ctx, charger.ID)
}

func (s *ChargerService) GetCharger(ctx context.Context, id uuid.UUID) (model.Charger, error) {
	return s.repo.FindByID(ctx, id)
}
