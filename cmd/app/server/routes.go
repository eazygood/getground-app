package server

import "github.com/gin-gonic/gin"

func initRoutes(router *gin.Engine, dependency *Dependecy) {
	router.POST("/guests", dependency.guestController.Create)
	router.GET("/guests/:id", dependency.guestController.GetById)
	router.PUT("/guests/:id", dependency.guestController.Update)
	router.DELETE("/guests/:id", dependency.guestController.Delete)

	router.POST("/guestlist/:guest_id", dependency.guestListController.Create)
	router.GET("/guestlist", dependency.guestListController.GetList)

	router.GET("/tables/empty_seats", dependency.tableController.GetEmptySeats)
}
