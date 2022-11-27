package server

import "github.com/gin-gonic/gin"

func initRoutes(router *gin.Engine, dependency *Dependecy) {
	router.POST("/guests", dependency.guestController.Create)
	router.PUT("/guests/:guest_id", dependency.guestController.Update)
	router.GET("/guests/:guest_id", dependency.guestController.GetById)
	router.GET("/guests/list", dependency.guestController.GetList)
	router.DELETE("/guests/:guest_id", dependency.guestController.Delete)

	router.POST("/guestlist/:guest_id", dependency.guestListController.Create)
	router.GET("/guestlist", dependency.guestListController.GetList)

	router.POST("/tables/", dependency.tableController.Create)
	router.PUT("/tables/:table_id", dependency.tableController.Update)
	router.DELETE("/tables/:table_id", dependency.tableController.Delete)
	router.GET("/tables/empty_seats", dependency.tableController.GetEmptySeats)
}
