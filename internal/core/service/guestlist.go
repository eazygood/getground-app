package service

import (
	"context"
	"fmt"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
)

type GuestListService struct {
	repository port.GuesListRepository
}

func NewGuestListService(repository port.GuesListRepository) port.GuestListService {
	return &GuestListService{
		repository: repository,
	}
}

func (g *GuestListService) FindAvailableTable(ctx context.Context, filter port.GetGuestListFilter) (*domain.Table, error) {
	table, err := g.repository.FindAvailableTable(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("get all tables: %w", err)
	}

	return table, nil
}

func (g *GuestListService) GetOccupiedSeats(ctx context.Context) ([]*domain.Table, error) {
	tables, err := g.repository.GetOccupiedSeats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all tables: %w", err)
	}

	return tables, nil
}
