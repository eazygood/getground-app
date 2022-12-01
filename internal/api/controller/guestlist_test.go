package controller

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/eazygood/getground-app/internal/api/controller/testutil"
	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	mockPort "github.com/eazygood/getground-app/mocks/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GuestListControllereSuite struct {
	suite.Suite
	*require.Assertions
	ctrl                 *gomock.Controller
	mockGuestService     *mockPort.MockGuestService
	mockTableService     *mockPort.MockTableService
	mockGuestListService *mockPort.MockGuestListService
	guestListController  GuestListController
}

func TestGuestListControllereSuite(t *testing.T) {
	suite.Run(t, new(GuestListControllereSuite))
}

func (g *GuestListControllereSuite) SetupTest() {
	g.Assertions = require.New(g.T())

	g.ctrl = gomock.NewController(g.T())
	g.mockGuestService = mockPort.NewMockGuestService(g.ctrl)
	g.mockTableService = mockPort.NewMockTableService(g.ctrl)
	g.mockGuestListService = mockPort.NewMockGuestListService(g.ctrl)
	g.guestListController = NewGuestListController(g.mockGuestService, g.mockTableService, g.mockGuestListService)
}

func (g *GuestListControllereSuite) TearDownTest() {
	g.ctrl.Finish()
}

func (g *GuestListControllereSuite) TestCreateGuestList() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	body := GuestListRequest{
		GuestID:            1,
		AccompanyingGuests: 5,
	}

	filter := port.GetGuestListFilter{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
	}

	testutil.MockJsonPost(c, body)

	guest := domain.Guest{
		ID:                 1,
		Name:               "Simon",
		AccompanyingGuests: 0,
		IsArrived:          false,
	}

	availableTable := domain.Table{
		ID:    2,
		Seats: 10,
	}

	g.mockGuestListService.EXPECT().FindAvailableTable(c, gomock.Eq(filter)).Return(&availableTable, nil).Times(1)
	g.mockGuestService.EXPECT().GetById(c, guest.ID).Return(&guest, nil).Times(1)
	g.mockGuestService.EXPECT().Update(c, guest.ID, &domain.Guest{AccompanyingGuests: uint16(body.AccompanyingGuests), IsArrived: true}).Return(nil).Times(1)
	g.mockTableService.EXPECT().Update(c, availableTable.ID, domain.Table{GuestID: &guest.ID}).Return(nil).Times(1)

	g.guestListController.Create(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	wantJson := `{"message":"success"}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *GuestListControllereSuite) TestCreateGuestListThrowErrorNoAvailableSeats() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	body := GuestListRequest{
		GuestID:            1,
		AccompanyingGuests: 1000,
	}

	filter := port.GetGuestListFilter{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
	}

	testutil.MockJsonPost(c, body)

	apiError := errors.New("no available seats")

	g.mockGuestListService.EXPECT().FindAvailableTable(c, gomock.Eq(filter)).Return(nil, apiError).Times(1)

	g.guestListController.Create(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusNotFound, w.Code)

	wantJson := `{"code":404,"message":"no available seats"}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *GuestListControllereSuite) TestCreateGuestListGuestAlreadHasTable() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	body := GuestListRequest{
		GuestID:            1,
		AccompanyingGuests: 5,
	}

	filter := port.GetGuestListFilter{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
	}

	testutil.MockJsonPost(c, body)

	guest := domain.Guest{
		ID:                 1,
		Name:               "Simon",
		AccompanyingGuests: 0,
		IsArrived:          true,
	}

	availableTable := domain.Table{
		ID:    2,
		Seats: 10,
	}

	g.mockGuestListService.EXPECT().FindAvailableTable(c, gomock.Eq(filter)).Return(&availableTable, nil).Times(1)
	g.mockGuestService.EXPECT().GetById(c, guest.ID).Return(&guest, nil).Times(1)

	g.guestListController.Create(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusBadRequest, w.Code)

	wantJson := `{"code":400,"message":"guest already has seats"}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *GuestListControllereSuite) TestGetListGuestList() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	testutil.MockJsonGet(c, []gin.Param{}, url.Values{})

	t := time.Now()
	guest1 := domain.Guest{
		ID:                 1,
		AccompanyingGuests: 10,
		TimeArrived:        &t,
		IsArrived:          true,
	}

	guest2 := domain.Guest{
		ID:                 2,
		AccompanyingGuests: 5,
		TimeArrived:        &t,
		IsArrived:          true,
	}

	occupiedSeats := []*domain.Table{
		{
			ID:      1,
			Seats:   10,
			GuestID: &guest1.ID,
		},
		{
			ID:      2,
			Seats:   10,
			GuestID: &guest2.ID,
		},
	}

	g.mockGuestListService.EXPECT().GetOccupiedSeats(c).Return(occupiedSeats, nil)

	g.guestListController.GetList(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	wantJson := `[{"id":1,"seats":10,"guest_id":1},{"id":2,"seats":10,"guest_id":2}]`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}
