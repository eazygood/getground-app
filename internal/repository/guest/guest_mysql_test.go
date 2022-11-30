package guest

import (
	"context"
	"database/sql"
	"regexp"
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

type GuestMysqlRepositorySuite struct {
	suite.Suite
	*require.Assertions
	DB                *gorm.DB
	mock              sqlmock.Sqlmock
	ctrl              *gomock.Controller
	mySqlGuestAdapter port.GuestRepository
}

func TestGuestMysqlRepositorySuite(t *testing.T) {
	suite.Run(t, new(GuestMysqlRepositorySuite))
}

func (g *GuestMysqlRepositorySuite) SetupTest() {
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
	g.mySqlGuestAdapter = NewMysqlGuestAdapter(g.DB)
}

func (u *GuestMysqlRepositorySuite) TearDownTest() {
	u.ctrl.Finish()
}

func (g *GuestMysqlRepositorySuite) TestCreateGuest() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	guest := &domain.Guest{
		Name:               "Tere",
		AccompanyingGuests: 0,
		TimeArrived:        nil,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "accompanying_guests", "time_arrived"}).AddRow(1, "Tere", 0, nil)
	g.mock.ExpectBegin()

	g.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `guests` (`name`,`accompanying_guests`,`time_arrived`,`is_arrived`) VALUES (?,?,?,?)")).
		WithArgs("Tere", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	g.mock.ExpectCommit()
	g.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE `guests`.`id` = ? ORDER BY `guests`.`id` LIMIT 1")).WithArgs(1).WillReturnRows(rows)

	_, err := g.mySqlGuestAdapter.Create(c, guest)

	require.NoError(g.T(), err)
}

func (g *GuestMysqlRepositorySuite) TestUpdateGuest() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	guest := &domain.Guest{
		Name:               "Tere",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
		IsArrived:          false,
	}

	g.mock.ExpectBegin()

	g.mock.ExpectExec("UPDATE `guests` SET (.+)  WHERE (.+)").
		WithArgs(guest.Name, guest.AccompanyingGuests, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	g.mock.ExpectCommit()

	err := g.mySqlGuestAdapter.Update(c, 1, guest)

	require.NoError(g.T(), err)
}

func (g *GuestMysqlRepositorySuite) TestDeleteGuest() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()
	id := 1

	g.mock.ExpectBegin()

	g.mock.ExpectExec(("DELETE FROM `guests` WHERE (.+)")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	g.mock.ExpectCommit()

	err := g.mySqlGuestAdapter.Delete(c, int64(id))

	require.NoError(g.T(), err)
}

func (g *GuestMysqlRepositorySuite) TestGetListUser() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	filters := port.GetGuestFilter{
		TimeArrived: true,
	}

	now := time.Now()

	expected := []*domain.Guest{
		{
			ID:                 1,
			Name:               "Tere",
			AccompanyingGuests: 10,
			TimeArrived:        &now,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "accompanying_guests", "time_arrived"}).AddRow(1, "Tere", 10, now)
	g.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE time_arrived IS NOT NULL")).WillReturnRows(rows)

	guests, err := g.mySqlGuestAdapter.GetAll(c, filters)

	require.NoError(g.T(), err)

	g.EqualValues(expected, guests)
}

func (g *GuestMysqlRepositorySuite) TestGetListWithOutFilterUser() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	filters := port.GetGuestFilter{
		TimeArrived: true,
	}

	now := time.Now()

	expected := []*domain.Guest{
		{
			ID:                 1,
			Name:               "Tere",
			AccompanyingGuests: 10,
			TimeArrived:        &now,
		},
		{
			ID:                 2,
			Name:               "Tere2",
			AccompanyingGuests: 0,
			TimeArrived:        nil,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "accompanying_guests", "time_arrived"}).
		AddRow(1, "Tere", 10, now).AddRow(2, "Tere2", 0, nil)
	g.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE time_arrived IS NOT NULL")).WillReturnRows(rows)

	guests, err := g.mySqlGuestAdapter.GetAll(c, filters)

	require.NoError(g.T(), err)

	g.EqualValues(expected, guests)
}

func (g *GuestMysqlRepositorySuite) TestGetById() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	expected := &domain.Guest{
		ID:                 1,
		Name:               "Tere",
		AccompanyingGuests: 0,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "accompanying_guests", "time_arrived"}).AddRow(1, "Tere", 0, nil)
	g.mock.ExpectQuery("^SELECT (.+) WHERE (.+)").WillReturnRows(rows)

	actual, err := g.mySqlGuestAdapter.GetById(c, 1)

	require.NoError(g.T(), err)

	g.EqualValues(expected, actual)
}
