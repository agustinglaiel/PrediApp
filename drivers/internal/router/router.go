package router

import (
	drivers "drivers/internal/api/drivers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, driverController *drivers.DriverController, driverEventController *drivers.DriverEventController) {
	// Use CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas relacionadas con drivers
	engine.POST("/drivers", driverController.CreateDriver)
	engine.GET("/drivers/:id", driverController.GetDriverByID)
	engine.PUT("/drivers/:id", driverController.UpdateDriver)
	engine.DELETE("/drivers/:id", driverController.DeleteDriver)
	engine.GET("/drivers", driverController.ListDrivers)
	engine.GET("/drivers/team/:teamName", driverController.ListDriversByTeam)
	engine.GET("/drivers/country/:countryCode", driverController.ListDriversByCountry)
	engine.GET("/drivers/fullname/:fullName", driverController.ListDriversByFullName)
	engine.GET("/drivers/acronym/:acronym", driverController.ListDriversByAcronym)

	// Rutas relacionadas con drivers_event
	engine.POST("/drivers-event", driverEventController.AddDriverToEvent)
	engine.DELETE("/drivers-event/:id", driverEventController.RemoveDriverFromEvent)
	engine.GET("/drivers-event/event/:event_id", driverEventController.ListDriversByEvent)
	engine.GET("/drivers-event/driver/:driver_id", driverEventController.ListEventsByDriver)

	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
