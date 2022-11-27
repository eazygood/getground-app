package port

import (
	"context"

	"github.com/eazygood/getground-app/internal/core/domain"
)

type GuestRepository interface {
	GetById(ctx context.Context, id int64) (*domain.Guest, error)
	GetAll(ctx context.Context) ([]*domain.Guest, error)
	Create(ctx context.Context, guest *domain.Guest) (*domain.Guest, error)
	Update(ctx context.Context, id int64, guest *domain.Guest) error
	Delete(ctx context.Context, id int64) error
}

type TableRepository interface {
	GetById(ctx context.Context, id int64) (*domain.Table, error)
	GetEmptySeats(ctx context.Context) ([]*domain.Table, error)
	Create(ctx context.Context, table *domain.Table) (*domain.Table, error)
	Update(ctx context.Context, id int64, table domain.Table) error
	Delete(ctx context.Context, id int64) error
}

type GuesListRepository interface {
	GetAll(ctx context.Context) ([]domain.GuestList, error)
	Create(ctx context.Context, guest domain.Guest, table domain.Table) error
}
