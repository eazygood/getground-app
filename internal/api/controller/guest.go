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

type GuestController interface {
	Create(request *gin.Context)
	Update(request *gin.Context)
	Delete(request *gin.Context)
	GetById(request *gin.Context)
	GetList(request *gin.Context)
}

type GuestRequest struct {
	Name               string `json:"name"`
	AccompanyingGuests uint16 `json:"accompanying_guests"`
	TimeArrived        string `json:"time_arrived,omitempty"`
}

func createFromCreateUpdateRequest(req GuestRequest) (*domain.Guest, error) {
	guest := domain.Guest{
		Name:               req.Name,
		AccompanyingGuests: req.AccompanyingGuests,
	}

	if req.TimeArrived != "" {
		t, err := strToTimePtr(req.TimeArrived)
		if err != nil {

			return nil, fmt.Errorf("invalid time arrived input")
		}

		guest.TimeArrived = t
	}

	return &guest, nil
}

type guestController struct {
	guestService port.GuestService
}

func NewGuestController(service port.GuestService) GuestController {
	return &guestController{
		guestService: service,
	}
}

func (c *guestController) Create(request *gin.Context) {
	body := GuestRequest{}
	if err := request.ShouldBindJSON(&body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	g, err := createFromCreateUpdateRequest(body)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	guest, err := c.guestService.Create(request, g)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusCreated, guest)
}

func (c *guestController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("guest_id"))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	body := GuestRequest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	g, err := createFromCreateUpdateRequest(body)
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	err = c.guestService.Update(ctx, int64(id), g)
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *guestController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("guest_id"))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	err = c.guestService.Delete(ctx, int64(id))
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *guestController) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("guest_id"))

	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	guest, err := c.guestService.GetById(ctx, int64(id))
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, guest)
}

func (c *guestController) GetList(ctx *gin.Context) {
	filters := port.GetGuestFilter{}

	if _, ok := ctx.GetQuery("arrived"); ok {
		filters.IsArrived = true
	}

	guests, err := c.guestService.GetList(ctx, filters)
	if err != nil {
		logAndAbort(ctx, errors.NewApiError(errors.Internal, err))
		return
	}

	ctx.JSON(http.StatusOK, guests)
}
