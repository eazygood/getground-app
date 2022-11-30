package service

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	mockPort "github.com/eazygood/getground-app/mocks/core/port"
	ports "github.com/eazygood/getground-app/mocks/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GuestListServiceSuite struct {
	suite.Suite
	*require.Assertions
	ctrl                    *gomock.Controller
	mockGuestListRepository *ports.MockGuesListRepository
	guestListService        port.GuestListService
}

func TestGuestListServiceSuite(t *testing.T) {
	suite.Run(t, new(GuestListServiceSuite))
}

func (g *GuestListServiceSuite) SetupTest() {
	g.Assertions = require.New(g.T())

	g.ctrl = gomock.NewController(g.T())
	g.mockGuestListRepository = mockPort.NewMockGuesListRepository(g.ctrl)

	g.guestListService = NewGuestListService(g.mockGuestListRepository)
}

func (g *GuestListServiceSuite) TearDownTest() {
	g.ctrl.Finish()
}

func (g *GuestListServiceSuite) TestGuestListFindAvailableTable() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	table := &domain.Table{
		Seats: 15,
	}

	g.mockGuestListRepository.EXPECT().FindAvailableTable(c, port.GetGuestListFilter{}).Return(table, nil).Times(1)

	actual, err := g.guestListService.FindAvailableTable(c, port.GetGuestListFilter{})
	g.NoError(err)

	g.EqualValues(table, actual)
}

func (g *GuestListServiceSuite) TestGuestListFindAvailableTableThrowError() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	err := errors.New("Mock Repository Error")
	g.mockGuestListRepository.EXPECT().FindAvailableTable(c, port.GetGuestListFilter{}).Return(nil, err).Times(1)

	_, err = g.guestListService.FindAvailableTable(c, port.GetGuestListFilter{})

	g.ErrorContains(err, "Mock Repository Error")
}

func (g *GuestListServiceSuite) TestGuestListGetOccupiedSeats() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	t := time.Now().UTC()

	guestJohn := &domain.Guest{
		ID:                 1,
		Name:               "John",
		AccompanyingGuests: 15,
		TimeArrived:        &t,
	}

	guestSimon := &domain.Guest{
		ID:                 1,
		Name:               "Simon",
		AccompanyingGuests: 15,
		TimeArrived:        &t,
	}

	tables := []*domain.Table{
		{
			Seats:   15,
			GuestID: &guestJohn.ID,
		},
		{
			Seats:   2,
			GuestID: &guestSimon.ID,
		},
	}

	g.mockGuestListRepository.EXPECT().GetOccupiedSeats(c).Return(tables, nil).Times(1)

	actual, err := g.guestListService.GetOccupiedSeats(c)
	g.NoError(err)

	g.EqualValues(tables, actual)
}

func (g *GuestListServiceSuite) TestGuestListGetOccupiedSeatsThrowError() {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	err := errors.New("Mock Repository Error")
	g.mockGuestListRepository.EXPECT().GetOccupiedSeats(c).Return(nil, err).Times(1)

	_, err = g.guestListService.GetOccupiedSeats(c)

	g.ErrorContains(err, "Mock Repository Error")
}
