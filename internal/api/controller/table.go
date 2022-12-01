package controller

import (
	"fmt"
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

type EmptySeatsResponse struct {
	EmptySeats int64 `json:"empty_seats"`
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

func (t *tableController) Create(ctx *gin.Context) {
	body := &TableCreateRequest{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	table, err := t.tableService.Create(ctx, &domain.Table{
		Seats: body.Seats,
	})

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusCreated, table)
}

func (t *tableController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("table_id"))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	body := TableUpdateeRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	tbl := domain.Table{}
	tbl.Seats = body.Seats
	tbl.GuestID = &body.GuestID

	if body.GuestID != 0 {
		table, err := t.tableService.GetById(ctx, int64(id))

		if err != nil {
			logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
			return
		}

		if table.GuestID != nil {
			logAndAbort(ctx, errors.NewApiError(errors.Internal, fmt.Errorf("table already has guest")))
			return
		}

		g, err := t.guestService.GetById(ctx, body.GuestID)

		if err != nil {
			logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
			return
		}

		if g.AccompanyingGuests > table.Seats {
			logAndAbort(ctx, errors.NewApiError(errors.Internal, fmt.Errorf("guest accompanying guests exceeded table available seats")))
			return
		}
	}

	err = t.tableService.Update(ctx, int64(id), tbl)
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (t *tableController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("table_id"))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	err = t.tableService.Delete(ctx, int64(id))
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (t *tableController) GetEmptySeats(ctx *gin.Context) {
	emptySeats, err := t.tableService.GetEmptySeats(ctx)
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, EmptySeatsResponse{EmptySeats: emptySeats})
}
