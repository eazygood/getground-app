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

type GuestServiceSuite struct {
	suite.Suite
	*require.Assertions
	ctrl                *gomock.Controller
	mockGuestRepository *ports.MockGuestRepository
	guestService        port.GuestService
}

func TestGuestServiceSuite(t *testing.T) {
	suite.Run(t, new(GuestServiceSuite))
}

func (g *GuestServiceSuite) SetupTest() {
	g.Assertions = require.New(g.T())

	g.ctrl = gomock.NewController(g.T())
	g.mockGuestRepository = mockPort.NewMockGuestRepository(g.ctrl)

	g.guestService = NewGuestService(g.mockGuestRepository)
}

func (g *GuestServiceSuite) TearDownTest() {
	g.ctrl.Finish()
}

func (g *GuestServiceSuite) TestCreateGuest() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	request := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
	}

	guest := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
	}

	g.mockGuestRepository.EXPECT().Create(c, gomock.Eq(request)).Return(guest, nil).Times(1)

	_, err := g.guestService.Create(c, guest)
	g.NoError(err)
}

func (g *GuestServiceSuite) TestCreateGuestThrowError() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	request := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
	}

	guest := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
	}

	err := errors.New("Mock Repository Error")
	g.mockGuestRepository.EXPECT().Create(c, gomock.Eq(request)).Return(guest, err).Times(1)

	_, err = g.guestService.Create(c, guest)
	g.ErrorContains(err, "Mock Repository Error")
}

func (g *GuestServiceSuite) TestGuestDelete() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	g.mockGuestRepository.EXPECT().Delete(c, int64(1)).Return(nil).Times(1)

	err := g.guestService.Delete(c, int64(1))
	g.NoError(err)
}

func (g *GuestServiceSuite) TestGuestGetById() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	guest := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
	}

	g.mockGuestRepository.EXPECT().GetById(c, int64(1)).Return(guest, nil).Times(1)

	actual, err := g.guestService.GetById(c, int64(1))
	g.NoError(err)

	g.EqualValues(guest, actual)
}

func (g *GuestServiceSuite) TestGuestGetByIdThrowError() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	guest := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
	}

	err := errors.New("Mock Repository Error")
	g.mockGuestRepository.EXPECT().GetById(c, int64(1)).Return(guest, err).Times(1)

	_, err = g.guestService.GetById(c, int64(1))
	g.ErrorContains(err, "Mock Repository Error")
}

func (g *GuestServiceSuite) TestGuestGetList() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	guests := []*domain.Guest{
		{
			Name:               "Simon",
			AccompanyingGuests: 10,
			TimeArrived:        nil,
		},
	}

	g.mockGuestRepository.EXPECT().GetAll(c, port.GetGuestFilter{}).Return(guests, nil).Times(1)

	actual, err := g.guestService.GetList(c, port.GetGuestFilter{})
	g.NoError(err)

	g.EqualValues(guests, actual)
}

func (g *GuestServiceSuite) TestGuestGetListThrowError() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	guests := []*domain.Guest{
		{
			Name:               "Simon",
			AccompanyingGuests: 10,
			TimeArrived:        nil,
		},
	}

	err := errors.New("Mock Repository Error")
	g.mockGuestRepository.EXPECT().GetAll(c, port.GetGuestFilter{}).Return(guests, err).Times(1)

	_, err = g.guestService.GetList(c, port.GetGuestFilter{})

	g.ErrorContains(err, "Mock Repository Error")
}

func (g *GuestServiceSuite) TestGuestUpdate() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	guest := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 10,
		TimeArrived:        nil,
	}

	g.mockGuestRepository.EXPECT().Update(c, int64(1), guest).Return(nil).Times(1)

	err := g.guestService.Update(c, int64(1), guest)
	g.NoError(err)
}
