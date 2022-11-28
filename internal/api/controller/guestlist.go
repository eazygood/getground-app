package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func (g *guestListController) Create(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("guest_id"))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	body := &GuestListRequest{}
	if err := request.ShouldBindJSON(&body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	filter := port.GetGuestListFilter{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	table, err := g.guestListService.FindAvailableTable(requestCtx, filter)

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	// get guest if it is registered(invited)
	guest, err := g.guestService.GetById(requestCtx, int64(id))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	if guest.TimeArrived != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, fmt.Errorf("guest already has seats")))
		return
	}

	// update guest with accompanying guest
	err = g.guestService.Update(requestCtx, guest.ID, &domain.Guest{
		AccompanyingGuests: uint16(body.AccompanyingGuests),
		TimeArrived:        toTimePtr(time.Now()),
	})

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	// update table with guest
	err = g.tableService.Update(requestCtx, table.ID, domain.Table{
		GuestID: &guest.ID,
	})

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (g *guestListController) GetList(request *gin.Context) {
	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	guestList, err := g.guestListService.GetOccupiedSeats(requestCtx)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, guestList)
}
