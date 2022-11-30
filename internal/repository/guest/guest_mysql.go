package guest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	v "github.com/eazygood/getground-app/internal/validator"
	"gorm.io/gorm"
)

type MysqlGuestAdapter struct {
	Conn *gorm.DB
}

func NewMysqlGuestAdapter(Conn *gorm.DB) port.GuestRepository {
	return &MysqlGuestAdapter{
		Conn: Conn,
	}
}

func (m *MysqlGuestAdapter) Create(ctx context.Context, guest *domain.Guest) (*domain.Guest, error) {
	if err := v.GetValidator().Struct(guest); err != nil {
		return nil, fmt.Errorf("failed to insert guest due to validation: %v", err)
	}

	err := m.Conn.Create(guest).Error

	if err != nil {
		return nil, fmt.Errorf("failed to insert guest: %v", err.Error())
	}

	g := &domain.Guest{}
	err = m.Conn.First(g, guest.ID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("record not found by id: %v", guest.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to delete guest by id (%v) %v", guest.ID, err.Error())
	}

	return guest, nil
}

func (m *MysqlGuestAdapter) Delete(ctx context.Context, id int64) error {
	err := m.Conn.Delete(&domain.Guest{}, id).Error

	if err != nil {
		return fmt.Errorf("failed to delete guest by id (%v) %v", id, err.Error())
	}

	return nil
}

func (m *MysqlGuestAdapter) GetById(ctx context.Context, id int64) (*domain.Guest, error) {
	guest := &domain.Guest{}
	err := m.Conn.First(guest, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("record not found by id: %v", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to delete guest by id (%v) %v", id, err.Error())
	}

	return guest, nil
}

func (m *MysqlGuestAdapter) Update(ctx context.Context, id int64, guest *domain.Guest) error {
	if guest.IsArrived {
		t := time.Now()
		guest.TimeArrived = &t
	}
	err := m.Conn.Model(&domain.Guest{}).Where("id = ?", id).Updates(guest).Error

	if err != nil {
		return fmt.Errorf("failed to update guest: %v", err.Error())
	}

	return nil
}

func (m *MysqlGuestAdapter) GetAll(ctx context.Context, filter port.GetGuestFilter) ([]*domain.Guest, error) {
	var guests []*domain.Guest

	conn := m.Conn

	if filter.TimeArrived {
		conn = conn.Where("time_arrived IS NOT NULL")
	}

	err := conn.Find(&guests).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get list of guests: %v", err.Error())
	}

	return guests, nil
}
