package port

import (
	"context"

	"github.com/eazygood/getground-app/internal/core/domain"
)

//go:generate mockgen -source service.go -destination=../../../mocks/core/port/service_mock.go -package ports
type GuestService interface {
	Create(ctx context.Context, g *domain.Guest) (*domain.Guest, error)
	Update(ctx context.Context, id int64, u *domain.Guest) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*domain.Guest, error)
	GetList(ctx context.Context, filter GetGuestFilter) ([]*domain.Guest, error)
}

type GuestListService interface {
	FindAvailableTable(ctx context.Context, filter GetGuestListFilter) (*domain.Table, error)
	GetOccupiedSeats(ctx context.Context) ([]*domain.Table, error)
}

type TableService interface {
	GetById(ctx context.Context, id int64) (*domain.Table, error)
	GetEmptySeats(ctx context.Context) ([]*domain.Table, error)
	Create(ctx context.Context, table *domain.Table) (*domain.Table, error)
	Update(ctx context.Context, id int64, table domain.Table) error
	Delete(ctx context.Context, id int64) error
}
