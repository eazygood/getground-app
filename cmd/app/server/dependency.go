package server

import (
	"github.com/eazygood/getground-app/internal/api/controller"
	"github.com/eazygood/getground-app/internal/config"
	"github.com/eazygood/getground-app/internal/core/service"
	mysql "github.com/eazygood/getground-app/internal/infrastructure/db"
	"github.com/eazygood/getground-app/internal/repository/guest"
	"github.com/eazygood/getground-app/internal/repository/guestlist"
	"github.com/eazygood/getground-app/internal/repository/table"
)

type Dependecy struct {
	guestController     controller.GuestController
	tableController     controller.TableController
	guestListController controller.GuestListController
}

func initDependencies(cfg *config.App) (*Dependecy, error) {
	// db connection
	db := mysql.InitDb(cfg)

	// repositories
	guestRepository := guest.NewMysqlGuestAdapter(db)
	tableRepository := table.NewMysqlTableAdapter(db)
	guestListRepository := guestlist.NewMysqlGuestListAdapter(db)

	// services
	guestService := service.NewGuestService(guestRepository)
	tableService := service.NewTableService(tableRepository)
	guestListService := service.NewGuestListService(guestListRepository)

	// controllers
	guestController := controller.NewGuestController(guestService)
	tableController := controller.NewTableController(tableService, guestService)
	guestLisController := controller.NewGuestListController(guestService, tableService, guestListService)

	return &Dependecy{
		guestController:     guestController,
		tableController:     tableController,
		guestListController: guestLisController,
	}, nil
}
