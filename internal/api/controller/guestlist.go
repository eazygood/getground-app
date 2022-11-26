package controller

import (
	"github.com/eazygood/getground-app/internal/core/port"
	"github.com/gin-gonic/gin"
)

type GuestListController interface {
	Create(request *gin.Context)
	GetList(request *gin.Context)
}

type guestListController struct {
	guestService port.GuestService
	tableService port.TableService
}

func NewGuestListController(guest port.GuestService, table port.TableService) GuestListController {
	return &guestListController{
		guestService: guest,
		tableService: table,
	}
}

// Create implements GuestListController
func (*guestListController) Create(request *gin.Context) {
	panic("unimplemented")
}

// GetList implements GuestListController
func (*guestListController) GetList(request *gin.Context) {
	panic("unimplemented")
}
