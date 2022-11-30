package service

import (
	"context"
	"fmt"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
)

type TableService struct {
	repository port.TableRepository
}

func NewTableService(repository port.TableRepository) port.TableService {
	return &TableService{
		repository: repository,
	}
}

func (srv *TableService) Create(ctx context.Context, table *domain.Table) (*domain.Table, error) {
	t, err := srv.repository.Create(ctx, table)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return t, nil
}

func (srv *TableService) Delete(ctx context.Context, id int64) error {
	if err := srv.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete table: %w", err)
	}

	return nil
}

func (srv *TableService) GetEmptySeats(ctx context.Context) ([]*domain.Table, error) {
	tables, err := srv.repository.GetEmptySeats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all tables: %w", err)
	}

	return tables, nil
}

func (srv *TableService) GetById(ctx context.Context, id int64) (*domain.Table, error) {
	table, err := srv.repository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get table: %w", err)
	}

	return table, nil
}

func (srv *TableService) Update(ctx context.Context, id int64, table domain.Table) error {
	if err := srv.repository.Update(ctx, id, table); err != nil {
		return fmt.Errorf("update table: %w", err)
	}

	return nil
}
