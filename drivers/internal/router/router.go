package router

import (
	drivers "prediapp.local/drivers/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, driverController *drivers.DriverController) {
	// Use CORS middleware
	// engine.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"}, // Incluye el gateway y el frontend
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	// Rutas relacionadas con drivers
	engine.POST("/drivers", driverController.CreateDriver)
	engine.GET("/drivers/:id", driverController.GetDriverByID)
	engine.PUT("/drivers/:id", driverController.UpdateDriver)
	engine.DELETE("/drivers/:id", driverController.DeleteDriver)
	engine.GET("/drivers", driverController.ListDrivers)
	engine.GET("/drivers/team", driverController.ListDriversByTeam)
	engine.GET("/drivers/country/:countryCode", driverController.ListDriversByCountry)
	engine.GET("/drivers/fullname/:fullName", driverController.ListDriversByFullName)
	engine.GET("/drivers/acronym/:acronym", driverController.ListDriversByAcronym)
	engine.GET("/drivers/external", driverController.FetchAllDriversFromExternalAPI)
	engine.GET("/drivers/number/:driver_number", driverController.GetDriverByNumber)

	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
