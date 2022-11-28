package table

import (
	"context"
	"errors"
	"fmt"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	v "github.com/eazygood/getground-app/internal/validator"
	"gorm.io/gorm"
)

type MysqlTableAdapter struct {
	Conn *gorm.DB
}

func NewMysqlTableAdapter(Conn *gorm.DB) port.TableRepository {
	return &MysqlTableAdapter{
		Conn: Conn,
	}
}

func (m *MysqlTableAdapter) Create(ctx context.Context, table *domain.Table) (*domain.Table, error) {
	if err := v.GetValidator().Struct(table); err != nil {
		return nil, fmt.Errorf("failed to insert guest due to validation: %v", err)
	}

	err := m.Conn.Create(table).Error

	if err != nil {
		return nil, fmt.Errorf("failed to table guest: %v", err.Error())
	}

	t := &domain.Table{}
	err = m.Conn.First(t, table.ID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("record not found by id: %v", table.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to delete guest by id (%v) %v", table.ID, err.Error())
	}

	return t, nil
}

func (m *MysqlTableAdapter) Update(ctx context.Context, id int64, table domain.Table) error {
	err := m.Conn.Model(&domain.Table{}).Where("id = ?", id).Updates(table).Error

	if err != nil {
		return fmt.Errorf("failed to update table: %v", err.Error())
	}

	return nil
}

func (m *MysqlTableAdapter) Delete(ctx context.Context, id int64) error {
	err := m.Conn.Delete(&domain.Table{}, id).Error

	if err != nil {
		return fmt.Errorf("failed to delete table by id (%v) %v", id, err.Error())
	}

	return nil
}

func (m *MysqlTableAdapter) GetEmptySeats(ctx context.Context) ([]*domain.Table, error) {
	var tables []*domain.Table
	err := m.Conn.Where("guest_id IS NULL").Find(&tables).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get list of tables: %v", err.Error())
	}

	return tables, nil
}

func (m *MysqlTableAdapter) GetOccupiedSeats(ctx context.Context) ([]*domain.Table, error) {
	var tables []*domain.Table

	err := m.Conn.Joins("JOIN guests ON guests.id = tables.guest_id").Where("tables.guest_id IS NOT NULL").Error

	if err != nil {
		return nil, fmt.Errorf("failed to get list of tables: %v", err.Error())
	}

	return tables, nil
}

func (m *MysqlTableAdapter) GetById(ctx context.Context, id int64) (*domain.Table, error) {
	table := &domain.Table{}
	err := m.Conn.First(table, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("record not found by id: %v", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to delete table by id (%v) %v", id, err.Error())
	}

	return table, nil
}
