package guestlist

import (
	"context"
	"errors"
	"fmt"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	"gorm.io/gorm"
)

type MysqlGuestListAdapter struct {
	Conn *gorm.DB
}

func NewMysqlGuestListAdapter(Conn *gorm.DB) port.GuesListRepository {
	return &MysqlGuestListAdapter{
		Conn: Conn,
	}
}

func (m *MysqlGuestListAdapter) FindAvailableTable(ctx context.Context, filter port.GetGuestListFilter) (*domain.Table, error) {
	table := domain.Table{}

	result := m.Conn.Debug().Where("seats >= ? AND guest_id IS NULL", filter.AccompanyingGuests).Find(&table)

	err := result.Error
	rows := result.RowsAffected

	if errors.Is(err, gorm.ErrRecordNotFound) || rows < 1 {
		return nil, fmt.Errorf("no available seats")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get available seats: %v", err.Error())
	}

	return &table, nil
}

func (m *MysqlGuestListAdapter) GetOccupiedSeats(ctx context.Context) ([]*domain.Table, error) {
	var tables []*domain.Table

	err := m.Conn.Debug().Preload("Guest").Joins("JOIN guests ON guests.id = guest_id").Where("guest_id IS NOT NULL").Find(&tables).Error

	// err := m.Conn.Debug().Joins("JOIN guests ON guests.id = guest_id").Where("guest_id IS NOT NULL").Find(guestList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get list of tables: %v", err.Error())
	}

	return tables, nil
}
