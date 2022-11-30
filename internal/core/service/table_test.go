package service

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	mockPort "github.com/eazygood/getground-app/mocks/core/port"
	ports "github.com/eazygood/getground-app/mocks/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TableServiceSuite struct {
	suite.Suite
	*require.Assertions
	ctrl                *gomock.Controller
	mockTableRepository *ports.MockTableRepository
	tableService        port.TableService
}

func TestTableServiceSuite(t *testing.T) {
	suite.Run(t, new(TableServiceSuite))
}

func (g *TableServiceSuite) SetupTest() {
	g.Assertions = require.New(g.T())

	g.ctrl = gomock.NewController(g.T())
	g.mockTableRepository = mockPort.NewMockTableRepository(g.ctrl)

	g.tableService = NewTableService(g.mockTableRepository)
}

func (g *TableServiceSuite) TearDownTest() {
	g.ctrl.Finish()
}

func (g *TableServiceSuite) TestCreateTable() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	request := &domain.Table{
		Seats: 10,
	}

	table := &domain.Table{
		Seats: 10,
	}

	g.mockTableRepository.EXPECT().Create(c, gomock.Eq(request)).Return(table, nil).Times(1)

	_, err := g.tableService.Create(c, table)
	g.NoError(err)
}

func (g *TableServiceSuite) TestCreateTableThrowError() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	request := &domain.Table{
		Seats: 10,
	}

	table := &domain.Table{
		Seats: 10,
	}

	err := errors.New("Mock Repository Error")

	g.mockTableRepository.EXPECT().Create(c, gomock.Eq(request)).Return(table, err).Times(1)

	_, err = g.tableService.Create(c, table)
	g.ErrorContains(err, "Mock Repository Error")
}

func (t *TableServiceSuite) TestTableDelete() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	t.mockTableRepository.EXPECT().Delete(c, int64(1)).Return(nil).Times(1)

	err := t.tableService.Delete(c, int64(1))
	t.NoError(err)
}

func (t *TableServiceSuite) TestTableGetById() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	table := &domain.Table{
		Seats: 15,
	}

	t.mockTableRepository.EXPECT().GetById(c, int64(1)).Return(table, nil).Times(1)

	actual, err := t.tableService.GetById(c, int64(1))
	t.NoError(err)

	t.EqualValues(table, actual)
}

func (t *TableServiceSuite) TestTableGetEmptySeats() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	tables := []*domain.Table{
		{
			Seats: 15,
		},
		{
			Seats: 2,
		},
	}

	t.mockTableRepository.EXPECT().GetEmptySeats(c).Return(tables, nil).Times(1)

	actual, err := t.tableService.GetEmptySeats(c)
	t.NoError(err)

	t.EqualValues(tables, actual)
}

func (t *TableServiceSuite) TestGuestUpdate() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	table := domain.Table{
		Seats: 15,
	}

	t.mockTableRepository.EXPECT().Update(c, int64(1), table).Return(nil).Times(1)

	err := t.tableService.Update(c, int64(1), table)
	t.NoError(err)
}
