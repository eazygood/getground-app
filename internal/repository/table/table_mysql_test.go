package table

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

type TableMysqlRepositorySuite struct {
	suite.Suite
	*require.Assertions
	DB                *gorm.DB
	mock              sqlmock.Sqlmock
	ctrl              *gomock.Controller
	mySqlTableAdapter port.TableRepository
}

func TestTableMysqlRepositorySuite(t *testing.T) {
	suite.Run(t, new(TableMysqlRepositorySuite))
}

func (t *TableMysqlRepositorySuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	t.Assertions = require.New(t.T())

	db, t.mock, err = sqlmock.New()
	require.NoError(t.T(), err)

	t.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})

	require.NoError(t.T(), err)

	t.ctrl = gomock.NewController(t.T())
	t.mySqlTableAdapter = NewMysqlTableAdapter(t.DB)
}

func (u *TableMysqlRepositorySuite) TearDownTest() {
	u.ctrl.Finish()
}

func (t *TableMysqlRepositorySuite) TestCreateTable() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	table := &domain.Table{
		Seats: 15,
	}

	rows := sqlmock.NewRows([]string{"id", "seat", "guest_id"}).AddRow(1, 15, nil)
	t.mock.ExpectBegin()

	t.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tables` (`seats`,`guest_id`) VALUES (?,?)")).
		WithArgs(15, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.mock.ExpectCommit()
	t.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ? ORDER BY `tables`.`id` LIMIT 1")).WithArgs(1).WillReturnRows(rows)

	_, err := t.mySqlTableAdapter.Create(c, table)

	require.NoError(t.T(), err)
}

func (t *TableMysqlRepositorySuite) TestCreateWithGuestIdTable() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	guest := &domain.Guest{
		ID:                 1,
		Name:               "Tere",
		AccompanyingGuests: 10,
	}
	table := &domain.Table{
		Seats:   15,
		GuestID: &guest.ID,
	}

	rows := sqlmock.NewRows([]string{"id", "seats", "guest_id"}).AddRow(1, 15, 1)

	t.mock.ExpectBegin()
	t.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tables` (`seats`,`guest_id`) VALUES (?,?)")).
		WithArgs(15, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.mock.ExpectCommit()
	t.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ? ORDER BY `tables`.`id` LIMIT 1")).WithArgs(1).WillReturnRows(rows)

	actual, err := t.mySqlTableAdapter.Create(c, table)

	require.NoError(t.T(), err)

	t.EqualValues(table, actual)
}

func (t *TableMysqlRepositorySuite) TestGetEmptySeats() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	expected := []*domain.Table{
		{
			ID:    1,
			Seats: 15,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "seats", "guest_id"}).AddRow(1, 15, nil)

	t.mock.ExpectQuery("^SELECT (.+) FROM `tables` WHERE (.+)").WillReturnRows(rows)

	actual, err := t.mySqlTableAdapter.GetEmptySeats(c)

	require.NoError(t.T(), err)

	t.EqualValues(expected, actual)
}

func (t *TableMysqlRepositorySuite) TestUpdateGuest() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	tableId := 1

	guest := domain.Guest{
		ID:                 77,
		Name:               "Tere",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
	}

	table := &domain.Table{
		Seats:   10,
		GuestID: &guest.ID,
	}

	t.mock.ExpectBegin()

	t.mock.ExpectExec("UPDATE `tables` SET (.+) WHERE (.+)").
		WithArgs(table.Seats, guest.ID, tableId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	t.mock.ExpectCommit()

	err := t.mySqlTableAdapter.Update(c, int64(tableId), *table)

	require.NoError(t.T(), err)
}

func (t *TableMysqlRepositorySuite) TestDeleteGuest() {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(1000))
	defer cancel()

	tableId := 1

	t.mock.ExpectBegin()

	t.mock.ExpectExec(("DELETE FROM `tables` WHERE (.+)")).
		WithArgs(tableId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	t.mock.ExpectCommit()

	err := t.mySqlTableAdapter.Delete(c, int64(tableId))

	require.NoError(t.T(), err)
}
