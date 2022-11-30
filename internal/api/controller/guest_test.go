package controller

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/eazygood/getground-app/internal/api/controller/testutil"
	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	mockPort "github.com/eazygood/getground-app/mocks/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GuestControllereSuite struct {
	suite.Suite
	*require.Assertions
	ctrl             *gomock.Controller
	mockGuestService *mockPort.MockGuestService
	guestController  GuestController
}

type HttpJsonResponse struct {
	Message string
}

func TestGuestControllereSuite(t *testing.T) {
	suite.Run(t, new(GuestControllereSuite))
}

func (g *GuestControllereSuite) SetupTest() {
	g.Assertions = require.New(g.T())

	g.ctrl = gomock.NewController(g.T())
	g.mockGuestService = mockPort.NewMockGuestService(g.ctrl)
	g.guestController = NewGuestController(g.mockGuestService)
}

func (g *GuestControllereSuite) TestCreateGuest() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	body := GuestRequest{
		Name: "Simon",
	}

	testutil.MockJsonPost(c, body)

	guestServiceData := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 0,
	}

	want := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 0,
	}

	g.mockGuestService.EXPECT().Create(c, gomock.Eq(guestServiceData)).Return(want, nil).Times(1)

	g.guestController.Create(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusCreated, w.Code)

	wantJson := `{"id":0,"name":"Simon","accompanying_guests":0,"time_arrived":null,"is_arrived":false}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *GuestControllereSuite) TearDownTest() {
	g.ctrl.Finish()
}

func (g *GuestControllereSuite) TestUpdateGuest() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "guest_id",
			Value: "1",
		},
	}

	body := GuestRequest{
		Name:               "Simon",
		AccompanyingGuests: 9999,
	}

	testutil.MockJsonPut(c, body, params)

	guestID := 1
	guestServiceData := &domain.Guest{
		Name:               "Simon",
		AccompanyingGuests: 9999,
	}

	g.mockGuestService.EXPECT().Update(c, int64(guestID), guestServiceData).Return(nil).Times(1)

	g.guestController.Update(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	got, err := io.ReadAll(res.Body)

	g.NoError(err)
	g.Equal(`{"message":"success"}`, string(got))
}

func (g *GuestControllereSuite) TestDeleteGuest() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "guest_id",
			Value: "1",
		},
	}

	testutil.MockJsonDelete(c, params)

	guestID := 1

	g.mockGuestService.EXPECT().Delete(c, int64(guestID)).Return(nil).Times(1)
	g.guestController.Delete(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	got, err := io.ReadAll(res.Body)

	g.NoError(err)
	g.Equal(`{"message":"success"}`, string(got))
}

func (g *GuestControllereSuite) TestGetByGuest() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "guest_id",
			Value: "1",
		},
	}

	testutil.MockJsonGet(c, params, url.Values{})

	guestID := 1
	guestServiceData := &domain.Guest{
		ID:                 1,
		Name:               "Simon",
		AccompanyingGuests: 20,
	}

	g.mockGuestService.EXPECT().GetById(c, int64(guestID)).Return(guestServiceData, nil).Times(1)
	g.guestController.GetById(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	got, err := io.ReadAll(res.Body)

	g.NoError(err)

	wantJson := `{"id":1,"name":"Simon","accompanying_guests":20,"time_arrived":null,"is_arrived":false}`
	g.Equal(wantJson, string(got))
}

func (g *GuestControllereSuite) TestGetListGuest() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "guest_id",
			Value: "1",
		},
	}

	testutil.MockJsonGet(c, params, url.Values{})

	guestServiceData := []*domain.Guest{
		{
			ID:                 1,
			Name:               "Simon",
			AccompanyingGuests: 20,
		},
		{
			ID:                 2,
			Name:               "John",
			AccompanyingGuests: 20,
		},
	}

	filter := port.GetGuestFilter{}

	g.mockGuestService.EXPECT().GetList(c, filter).Return(guestServiceData, nil).Times(1)
	g.guestController.GetList(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	got, err := io.ReadAll(res.Body)

	g.NoError(err)

	wantJson := `[{"id":1,"name":"Simon","accompanying_guests":20,"time_arrived":null,"is_arrived":false},{"id":2,"name":"John","accompanying_guests":20,"time_arrived":null,"is_arrived":false}]`
	g.Equal(wantJson, string(got))
}
