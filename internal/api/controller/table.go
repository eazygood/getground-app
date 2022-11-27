package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/eazygood/getground-app/internal/core/domain"
	"github.com/eazygood/getground-app/internal/core/port"
	"github.com/eazygood/getground-app/internal/errors"
	"github.com/gin-gonic/gin"
)

type TableController interface {
	Create(request *gin.Context)
	GetEmptySeats(request *gin.Context)
	Update(request *gin.Context)
	Delete(request *gin.Context)
}

type TableRequest struct {
	Seats   uint16 `json:"seats"`
	GuestID int64  `json:"guest_id"`
}

type tableController struct {
	tableService port.TableService
	guestService port.GuestService
}

func NewTableController(tableService port.TableService, guestService port.GuestService) TableController {
	return &tableController{
		tableService: tableService,
		guestService: guestService,
	}
}

func (t *tableController) Create(request *gin.Context) {
	// body := &TableRequest{}
	body := &domain.Table{}
	if err := request.ShouldBindJSON(body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	// tbl := domain.Table{}
	// tbl.Seats = body.Seats

	// if body.GuestID != 0 {
	// 	g, err := t.guestService.GetById(requestCtx, body.GuestID)

	// 	if err != nil {
	// 		logAndAbort(request, errors.NewApiError(errors.Internal, err))
	// 		return
	// 	}

	// 	tbl.Guest = g
	// }

	table, err := t.tableService.Create(requestCtx, body)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusCreated, table)
}

func (t *tableController) Update(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("table_id"))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	body := TableRequest{}
	if err := request.ShouldBind(&body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	tbl := domain.Table{}
	tbl.Seats = body.Seats

	// check is table already has guest_id, if not return error
	// check guest_id if exist, if not return error

	// if body.GuestID != 0 {
	// 	g, err := t.guestService.GetById(requestCtx, body.GuestID)

	// 	if err != nil {
	// 		logAndAbort(request, errors.NewApiError(errors.Internal, err))
	// 		return
	// 	}

	// 	tbl.Guest = g
	// }

	err = t.tableService.Update(requestCtx, int64(id), tbl)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (t *tableController) Delete(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("table_id"))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	err = t.tableService.Delete(requestCtx, int64(id))
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (t *tableController) GetEmptySeats(request *gin.Context) {
	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	emptySeats, err := t.tableService.GetEmptySeats(requestCtx)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, emptySeats)
}
