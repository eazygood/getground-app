package service

import (
	"context"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
)

type guestListService struct {
	repository port.GuesListRepository
}

func NewGuestListService(repository port.GuesListRepository) port.GuesListRepository {
	return &guestListService{
		repository: repository,
	}
}

// Create implements port.GuesListRepository
func (*guestListService) Create(ctx context.Context, guest domain.Guest, table domain.Table) error {
	panic("unimplemented")
}

// GetAll implements port.GuesListRepository
func (*guestListService) GetAll(ctx context.Context) ([]domain.GuestList, error) {
	panic("unimplemented")
}
