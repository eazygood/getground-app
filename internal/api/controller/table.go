package controller

import (
	"github.com/eazygood/getground-app/internal/core/port"
	"github.com/gin-gonic/gin"
)

type TableController interface {
	Create(request *gin.Context)
	GetEmptySeats(request *gin.Context)
	Update(request *gin.Context)
	Delete(request *gin.Context)
}

type tableController struct {
	tableService port.TableService
}

func NewTableController(service port.TableService) TableController {
	return &tableController{
		tableService: service,
	}
}

// Create implements TableController
func (*tableController) Create(request *gin.Context) {
	panic("unimplemented")
}

// Delete implements TableController
func (*tableController) Delete(request *gin.Context) {
	panic("unimplemented")
}

// Update implements TableController
func (*tableController) Update(request *gin.Context) {
	panic("unimplemented")
}

func (*tableController) GetEmptySeats(request *gin.Context) {
	panic("unimplemented")
}
