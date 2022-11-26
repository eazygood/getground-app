package controller

import (
	"github.com/eazygood/getground-app/internal/core/port"
	"github.com/gin-gonic/gin"
)

type GuestController interface {
	Create(request *gin.Context)
	Update(request *gin.Context)
	Delete(request *gin.Context)
	GetById(request *gin.Context)
	GetList(request *gin.Context)
}

type guestController struct {
	guestService port.GuestService
}

func NewGuestController(service port.GuestService) GuestController {
	return &guestController{
		guestService: service,
	}
}

// Create implements GuestController
func (*guestController) Create(request *gin.Context) {
	panic("unimplemented")
}

// Delete implements GuestController
func (*guestController) Delete(request *gin.Context) {
	panic("unimplemented")
}

// GetById implements GuestController
func (*guestController) GetById(request *gin.Context) {
	panic("unimplemented")
}

// GetList implements GuestController
func (*guestController) GetList(request *gin.Context) {
	panic("unimplemented")
}

// Update implements GuestController
func (*guestController) Update(request *gin.Context) {
	panic("unimplemented")
}
