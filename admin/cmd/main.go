package main

import (
	drC "admin/internal/api/drivers"
	prC "admin/internal/api/prodes"
	seC "admin/internal/api/sessions"
	drR "admin/internal/repository/drivers"
	prR "admin/internal/repository/prodes"
	seR "admin/internal/repository/sessions"
	"admin/internal/router"
	drS "admin/internal/service/drivers"
	prS "admin/internal/service/prodes"
	seS "admin/internal/service/sessions"
	"admin/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)



func main() {
	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	// Iniciar el motor de la base de datos para migrar tablas
	utils.StartDbEngine()

	// Crear repositorios
	sessionRepo := seR.NewSessionRepository(db)
	driverRepo := drR.NewDriverRepository(db)
	driverEventRepo := drR.NewDriverEventRepository(db)
	prodeRepo := prR.NewProdeRepository(db)

	// Crear servicios
	sessionService := seS.NewSessionService(sessionRepo)
	driverService := drS.NewDriverService(driverRepo)
	driverEventService := drS.NewDriverEventService(driverEventRepo, driverRepo)
	prodeService := prS.NewPrediService(prodeRepo, sessionService)

	// Crear controladores
	sessionController := seC.NewSessionController(sessionService)
	driverController := drC.NewDriverController(driverService)
	driverEventController := drC.NewDriverEventController(driverEventService)
	prodeController := prC.NewProdeController(prodeService)

	// Configurar el router
	engine := gin.Default()
	router.MapUrls(engine, sessionController, driverController, driverEventController, prodeController)

	// Ejecutar el servidor
	if err := engine.Run(":8080"); err != nil {
		panic(err)
	}
}
