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

type GuestController interface {
	Create(request *gin.Context)
	Update(request *gin.Context)
	Delete(request *gin.Context)
	GetById(request *gin.Context)
	GetList(request *gin.Context)
}

type CreateGuestRequest struct {
	FirstName          string `json:"name"`
	AccompanyingGuests uint16 `json:"accompanying_guests"`
	TimeArrived        uint16 `json:"time_arrived,omitempty"`
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
	body := domain.Guest{}
	if err := request.BindJSON(&body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	guest, err := c.guestService.Create(requestCtx, &body)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusCreated, guest)
}

func (c *guestController) Delete(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("guest_id"))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	err = c.guestService.Delete(requestCtx, int64(id))
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *guestController) GetById(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("guest_id"))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	guest, err := c.guestService.GetById(requestCtx, int64(id))
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, guest)
}

func (c *guestController) GetList(request *gin.Context) {
	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	guests, err := c.guestService.GetList(requestCtx)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, guests)
}

func (c *guestController) Update(request *gin.Context) {
	id, err := strconv.Atoi(request.Param("guest_id"))

	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	body := domain.Guest{}
	if err := request.BindJSON(&body); err != nil {
		logAndAbort(request, errors.NewApiError(errors.InvalidInput, err))
		return
	}

	requestCtx, cancel := context.WithTimeout(request, requestTimeout)
	defer cancel()

	err = c.guestService.Update(requestCtx, int64(id), &body)
	if err != nil {
		logAndAbort(request, errors.NewApiError(errors.Internal, err))
		return
	}

	request.JSON(http.StatusOK, gin.H{"message": "success"})
}
