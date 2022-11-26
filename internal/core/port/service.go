package port

import (
	"context"

	"github.com/eazygood/getground-app/internal/core/domain"
)

type GuestService interface {
	Create(ctx context.Context, g *domain.Guest) (*domain.Guest, error)
	Update(ctx context.Context, id int64, u *domain.Guest) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*domain.Guest, error)
	GetList(ctx context.Context) ([]*domain.Guest, error)
}

type GuestListService interface {
	GetAll(ctx context.Context) ([]domain.GuestList, error)
	Create(ctx context.Context, guest domain.Guest, table domain.Table) error
}

type TableService interface {
	GetById(ctx context.Context, id int64) (*domain.Table, error)
	GetAll(ctx context.Context) ([]*domain.Table, error)
	Create(ctx context.Context, table domain.Table) error
	Update(ctx context.Context, id int64, table domain.Table) error
	Delete(ctx context.Context, id int64) error
}
