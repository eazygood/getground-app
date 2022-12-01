package controller

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/eazygood/getground-app/internal/api/controller/testutil"
	"github.com/eazygood/getground-app/internal/core/domain"
	mockPort "github.com/eazygood/getground-app/mocks/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TableControllereSuite struct {
	suite.Suite
	*require.Assertions
	ctrl             *gomock.Controller
	mockTableService *mockPort.MockTableService
	mockGuestService *mockPort.MockGuestService
	tableController  TableController
}

func TestTabletControllereSuite(t *testing.T) {
	suite.Run(t, new(TableControllereSuite))
}

func (g *TableControllereSuite) SetupTest() {
	g.Assertions = require.New(g.T())
	g.ctrl = gomock.NewController(g.T())
	g.mockTableService = mockPort.NewMockTableService(g.ctrl)
	g.mockGuestService = mockPort.NewMockGuestService(g.ctrl)
	g.tableController = NewTableController(g.mockTableService, g.mockGuestService)
}

func (g *TableControllereSuite) TearDownTest() {
	g.ctrl.Finish()
}

func (g *TableControllereSuite) TestCreateTable() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	body := TableCreateRequest{
		Seats: 15,
	}

	testutil.MockJsonPost(c, body)

	tableData := &domain.Table{
		Seats: 15,
	}

	want := &domain.Table{
		ID:    1,
		Seats: 15,
	}
	g.mockTableService.EXPECT().Create(c, gomock.Eq(tableData)).Return(want, nil).Times(1)

	g.tableController.Create(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusCreated, w.Code)

	wantJson := `{"id":1,"seats":15,"guest_id":null}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *TableControllereSuite) TestUpdateTable() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "table_id",
			Value: "1",
		},
	}

	body := TableUpdateeRequest{
		Seats: 15,
	}

	testutil.MockJsonPut(c, body, params)

	tableId := 1
	guest := domain.Guest{}
	tableData := domain.Table{
		Seats:   15,
		GuestID: &guest.ID,
	}

	g.mockTableService.EXPECT().Update(c, int64(tableId), gomock.Eq(tableData)).Return(nil).Times(1)

	g.tableController.Update(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	wantJson := `{"message":"success"}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *TableControllereSuite) TestUpdateTableWithGuestId() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "table_id",
			Value: "1",
		},
	}

	body := TableUpdateeRequest{
		Seats:   15,
		GuestID: 1,
	}

	testutil.MockJsonPut(c, body, params)

	tableId := 1
	guest := domain.Guest{
		ID:                 1,
		AccompanyingGuests: 0,
		IsArrived:          false,
	}
	tableData := domain.Table{
		Seats:   15,
		GuestID: &guest.ID,
	}

	returnTableData := domain.Table{
		Seats: 15,
	}

	g.mockTableService.EXPECT().GetById(c, int64(tableId)).Return(&returnTableData, nil).Times(1)
	g.mockGuestService.EXPECT().GetById(c, int64(tableId)).Return(&guest, nil).Times(1)
	g.mockTableService.EXPECT().Update(c, int64(tableId), gomock.Eq(tableData)).Return(nil).Times(1)

	g.tableController.Update(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	wantJson := `{"message":"success"}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *TableControllereSuite) TestUpdateTableWithGuestIdAndAccompanyingGuestExceeded() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "table_id",
			Value: "1",
		},
	}

	body := TableUpdateeRequest{
		Seats:   15,
		GuestID: 1,
	}

	testutil.MockJsonPut(c, body, params)

	tableId := 1
	guest := domain.Guest{
		ID:                 1,
		AccompanyingGuests: 1000,
		IsArrived:          false,
	}

	returnTableData := domain.Table{
		Seats: 15,
	}

	g.mockTableService.EXPECT().GetById(c, int64(tableId)).Return(&returnTableData, nil).Times(1)
	g.mockGuestService.EXPECT().GetById(c, int64(tableId)).Return(&guest, nil).Times(1)

	g.tableController.Update(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(500, w.Code)

	wantJson := `{"code":500,"message":"guest accompanying guests exceeded table available seats"}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}

func (g *TableControllereSuite) TestDeleteTable() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	params := []gin.Param{
		{
			Key:   "table_id",
			Value: "1",
		},
	}

	testutil.MockJsonDelete(c, params)

	tableID := 1

	g.mockTableService.EXPECT().Delete(c, int64(tableID)).Return(nil).Times(1)
	g.tableController.Delete(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	got, err := io.ReadAll(res.Body)

	g.NoError(err)
	g.Equal(`{"message":"success"}`, string(got))
}

func (g *TableControllereSuite) TestGetEmptySeatsTable() {
	w := httptest.NewRecorder()
	c := testutil.GetTestGinContext(w)

	testutil.MockJsonGet(c, []gin.Param{}, url.Values{})

	g.mockTableService.EXPECT().GetEmptySeats(c).Return(int64(17), nil)

	g.tableController.GetEmptySeats(c)

	res := w.Result()
	defer res.Body.Close()

	g.EqualValues(http.StatusOK, w.Code)

	wantJson := `{"empty_seats":17}`
	got, _ := io.ReadAll(res.Body)

	g.Equal(wantJson, string(got))
}
