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

type GuestListController interface {
	Create(request *gin.Context)
	GetList(request *gin.Context)
}

type GuestListRequest struct {
	TableID            int `json:"table_id"`
	AccompanyingGuests int `json:"accompanying_guests"`
}

type guestListController struct {
	guestService     port.GuestService
	tableService     port.TableService
	guestListService port.GuestListService
}

func NewGuestListController(guest port.GuestService, table port.TableService, guestList port.GuestListService) GuestListController {
	return &guestListController{
		guestService:     guest,
		tableService:     table,
		guestListService: guestList,
	}
}

func (g *guestListController) Create(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("guest_id"))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	body := &GuestListRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	filter := port.GetGuestListFilter{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
	}

	table, err := g.guestListService.FindAvailableTable(ctx, filter)

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	// get guest if it is registered(invited)
	guest, err := g.guestService.GetById(ctx, int64(id))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	if guest.IsArrived {
		logAndAbort(ctx, errors.NewApiError(errors.NotFound, fmt.Errorf("guest already has seats")))
		return
	}

	// update guest with accompanying guest
	err = g.guestService.Update(ctx, guest.ID, &domain.Guest{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
		IsArrived:          true,
	})

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	// update table with guest
	err = g.tableService.Update(ctx, table.ID, domain.Table{
		GuestID: &guest.ID,
	})

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (g *guestListController) GetList(ctx *gin.Context) {
	guestList, err := g.guestListService.GetOccupiedSeats(ctx)
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, guestList)
}
