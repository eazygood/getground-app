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

type TableCreateRequest struct {
	Seats uint16 `json:"seats"`
}
type TableUpdateeRequest struct {
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
	body := &TableCreateRequest{}
	if err := request.ShouldBindJSON(body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	table, err := t.tableService.Create(requestCtx, &domain.Table{
		Seats: body.Seats,
	})

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

	body := TableUpdateeRequest{}
	if err := request.ShouldBind(&body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	tbl := domain.Table{}
	tbl.Seats = body.Seats
	tbl.GuestID = &body.GuestID

	if body.GuestID != 0 {
		table, err := t.tableService.GetById(requestCtx, int64(id))

		if err != nil {
			logAndAbort(request, errors.NewApiError(errors.Internal, err))
			return
		}

		if table.GuestID != nil {
			logAndAbort(request, errors.NewApiError(errors.Internal, err))
			return
		}

		g, err := t.guestService.GetById(requestCtx, body.GuestID)

		if err != nil {
			logAndAbort(request, errors.NewApiError(errors.Internal, err))
			return
		}

		if g.AccompanyingGuests > table.Seats {
			logAndAbort(request, errors.NewApiError(errors.Internal, err))
			return
		}
	}

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
