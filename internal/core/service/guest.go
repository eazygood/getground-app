package service

import (
	"context"
	"fmt"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
)

type GuestService struct {
	repository port.GuestRepository
}

func NewGuestService(repository port.GuestRepository) port.GuestService {
	return &GuestService{repository: repository}
}

func (srv *GuestService) Create(ctx context.Context, guest *domain.Guest) (*domain.Guest, error) {
	guest, err := srv.repository.Create(ctx, guest)
	if err != nil {
		return nil, fmt.Errorf("create guest: %w", err)
	}

	return guest, nil
}

func (srv *GuestService) Delete(ctx context.Context, id int64) error {
	if err := srv.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete guest: %w", err)
	}

	return nil
}

func (srv *GuestService) GetById(ctx context.Context, id int64) (*domain.Guest, error) {
	guest, err := srv.repository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get guest: %w", err)
	}

	return guest, nil
}

func (srv *GuestService) GetList(ctx context.Context, filter port.GetGuestFilter) ([]*domain.Guest, error) {
	guests, err := srv.repository.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get all guests: %w", err)
	}

	return guests, nil
}

func (srv *GuestService) Update(ctx context.Context, id int64, guest *domain.Guest) error {
	if err := srv.repository.Update(ctx, id, guest); err != nil {
		return fmt.Errorf("update guest: %w", err)
	}

	return nil
}
