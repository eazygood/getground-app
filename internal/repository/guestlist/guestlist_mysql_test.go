package guestlist

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GuestListMysqlRepositorySuite struct {
	suite.Suite
	*require.Assertions
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	ctrl           *gomock.Controller
	mySqlGuestList port.GuesListRepository
}

func TestTableMysqlRepositorySuite(t *testing.T) {
	suite.Run(t, new(GuestListMysqlRepositorySuite))
}

func (g *GuestListMysqlRepositorySuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	g.Assertions = require.New(g.T())

	db, g.mock, err = sqlmock.New()
	require.NoError(g.T(), err)

	g.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})

	require.NoError(g.T(), err)

	g.ctrl = gomock.NewController(g.T())
	g.mySqlGuestList = NewMysqlGuestListAdapter(g.DB)
}

func (u *GuestListMysqlRepositorySuite) TearDownTest() {
	u.ctrl.Finish()
}

func (g *GuestListMysqlRepositorySuite) TestGetOccupiedSeats() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	guest := &domain.Guest{
		ID:                 1,
		Name:               "Tere",
		AccompanyingGuests: 10,
	}

	expected := []*domain.Table{
		{
			ID:      1,
			Seats:   15,
			GuestID: &guest.ID,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "seats", "guest_id"}).AddRow(1, 15, 1)

	g.mock.ExpectQuery("^SELECT (.+) FROM `tables` JOIN (.+) WHERE (.+)").WillReturnRows(rows)
	g.mock.ExpectQuery("^SELECT (.+) FROM `guests` WHERE (.+)").WillReturnRows(rows)

	actual, err := g.mySqlGuestList.GetOccupiedSeats(c)

	require.NoError(g.T(), err)

	g.EqualValues(expected, actual)
}

func (g *GuestListMysqlRepositorySuite) TestFindAvailableTable() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	expected := &domain.Table{
		ID:    1,
		Seats: 15,
	}

	filter := port.GetGuestListFilter{
		AccompanyingGuests: 10,
	}

	rows := sqlmock.NewRows([]string{"id", "seats", "guest_id"}).AddRow(1, 15, nil)

	g.mock.ExpectQuery("^SELECT (.+) FROM `tables` WHERE (.+)").WillReturnRows(rows)

	actual, err := g.mySqlGuestList.FindAvailableTable(c, filter)

	require.NoError(g.T(), err)

	g.EqualValues(expected, actual)
}

func (g *GuestListMysqlRepositorySuite) TestFindAvailableTableNotFound() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	rows := sqlmock.NewRows([]string{"id", "seats", "guest_id"})

	g.mock.ExpectQuery("^SELECT (.+) FROM `tables` WHERE (.+)").WillReturnRows(rows)

	_, err := g.mySqlGuestList.FindAvailableTable(c, port.GetGuestListFilter{})

	g.ErrorContains(err, "no available seats")
}
