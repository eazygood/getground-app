package port

import (
	"context"

	"github.com/eazygood/getground-app/internal/core/domain"
)

type GetGuestFilter struct {
	TimeArrived bool `json:"time_arrived"`
}

//go:generate mockgen -source repository.go -destination=../../../mocks/core/port/repository_mock.go -package ports
type GuestRepository interface {
	GetById(ctx context.Context, id int64) (*domain.Guest, error)
	GetAll(ctx context.Context, filter GetGuestFilter) ([]*domain.Guest, error)
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

type GetGuestListFilter struct {
	AccompanyingGuests uint16 `json:"accompanying_guests"`
}

type GuesListRepository interface {
	FindAvailableTable(ctx context.Context, filter GetGuestListFilter) (*domain.Table, error)
	GetOccupiedSeats(ctx context.Context) ([]*domain.Table, error)
}
