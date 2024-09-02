package router

import (
	drivers "admin/internal/api/drivers"
	prodes "admin/internal/api/prodes"
	api "admin/internal/api/sessions"
	users "admin/internal/api/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, sessionController *api.SessionController, driverController *drivers.DriverController, driverEventController *drivers.DriverEventController, prodeController *prodes.ProdeController, userController *users.UserController) {
	// Use CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas relacionadas con sesiones
	engine.POST("/sessions", sessionController.CreateSession)
	engine.GET("/sessions/:id", sessionController.GetSessionById)
	engine.PUT("/sessions/:id", sessionController.UpdateSessionById)
	engine.DELETE("/sessions/:id", sessionController.DeleteSessionById)
	engine.GET("/sessions/year/:year", sessionController.ListSessionsByYear)
	engine.GET("/sessions/circuit/:circuitKey", sessionController.ListSessionsByCircuitKey)
	engine.GET("/sessions/country/:countryCode", sessionController.ListSessionsByCountryCode)
	engine.GET("/sessions/upcoming", sessionController.ListUpcomingSessions)
	// engine.GET("/sessions/date-range", sessionController.ListSessionsBetweenDates) VER DESPUES
	engine.GET("/sessions/name-type", sessionController.FindSessionsByNameAndType)
	engine.GET("/sessions/:id/name-type", sessionController.GetSessionNameAndTypeById)
	engine.GET("/sessions", sessionController.GetAllSessions)

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

	// Rutas relacionadas con prodes
	engine.POST("/prodes/carrera", prodeController.CreateProdeCarrera)
	engine.POST("/prodes/session", prodeController.CreateProdeSession)
	engine.PUT("/prodes/carrera/:id", prodeController.UpdateProdeCarrera)
	engine.PUT("/prodes/session/:id", prodeController.UpdateProdeSession)
	engine.DELETE("/prodes/:id", prodeController.DeleteProdeById)
	engine.GET("/prodes/user/:user_id", prodeController.GetProdesByUserId)

	// Rutas relacionadas con usuarios
	engine.POST("/users/signup", userController.SignUp)
	engine.POST("/users/login", userController.Login)
	engine.POST("/users/oauth", userController.OAuthSignIn)
	engine.GET("/users/:id", userController.GetUserByID)
	engine.GET("/users/username/:username", userController.GetUserByUsername)
	engine.GET("/users", userController.GetUsers)
	engine.PUT("/users/:id", userController.UpdateUserByID)
	engine.PUT("/users/username/:username", userController.UpdateUserByUsername)
	engine.DELETE("/users/:id", userController.DeleteUserByID)
	engine.DELETE("/users/username/:username", userController.DeleteUserByUsername)
	engine.PUT("/users/:id/role", userController.UpdateUserRoleByID)
	engine.PUT("/users/:id/deactivate", userController.DeactivateUserByID)
	engine.PUT("/users/:id/reactivate", userController.ReactivateUserByID)


	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
