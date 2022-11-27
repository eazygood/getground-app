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
	// check if table already has guest
	// check available tables if seats is null
	// compare available seats with guest's accompanying guests
	// add guest to guest to table
	panic("unimplemented")
}

// GetList implements GuestListController
func (*guestListController) GetList(request *gin.Context) {
	panic("unimplemented")
}
