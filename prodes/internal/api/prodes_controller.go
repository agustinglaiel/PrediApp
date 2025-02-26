package api

import (
	"log"
	"net/http"
	dto "prodes/internal/dto"
	prodes "prodes/internal/service"
	e "prodes/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProdeController struct {
	prodeService prodes.ProdeServiceInterface
}

func NewProdeController(prodeService prodes.ProdeServiceInterface) *ProdeController {
	return &ProdeController{
		prodeService: prodeService,
	}
}

func (c *ProdeController) CreateProdeCarrera(ctx *gin.Context) {
	var request dto.CreateProdeCarreraDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	response, err := c.prodeService.CreateProdeCarrera(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *ProdeController) CreateProdeSession(ctx *gin.Context) {
    log.Printf("Creando pronóstico de sesión")
	var request dto.CreateProdeSessionDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
        log.Printf("Error en el bind del JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

    log.Printf("Data recibida: %+v", request)

	response, err := c.prodeService.CreateProdeSession(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *ProdeController) UpdateProdeCarrera(ctx *gin.Context) {
	// Capturar el prode_id desde la URL
	prodeIDParam := ctx.Param("prode_id")

	// Convertir el prode_id de string a int
	prodeID, err := strconv.Atoi(prodeIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid prode_id in URL"))
		return
	}

	// Parsear el JSON del request body a un DTO
	var request dto.UpdateProdeCarreraDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	// Asignar el prodeID capturado al DTO antes de enviarlo al servicio
	request.ProdeID = prodeID

	// Llamar al servicio para actualizar el ProdeCarrera
	response, err := c.prodeService.UpdateProdeCarrera(ctx.Request.Context(), request)
	if err != nil {
		if apiErr, ok := err.(e.ApiError); ok {
			ctx.JSON(apiErr.Status(), apiErr)
		} else {
			// Si no lo es, retornamos un error interno genérico
			ctx.JSON(http.StatusInternalServerError, e.NewInternalServerApiError("Unexpected error", err))
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) UpdateProdeSession(ctx *gin.Context) {
	// Obtener `session_id` desde los parámetros de la URL
	sessionID, err := strconv.Atoi(ctx.Param("session_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session_id parameter"))
		return
	}

	var request dto.UpdateProdeSessionDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	// Asignar `sessionID` obtenido de la URL al DTO
	request.SessionID = sessionID

	response, err := c.prodeService.UpdateProdeSession(ctx.Request.Context(), request)
	if err != nil {
		// Hacer type assertion para verificar si `err` es del tipo `ApiError`
		if apiErr, ok := err.(e.ApiError); ok {
			ctx.JSON(apiErr.Status(), apiErr)
		} else {
			ctx.JSON(http.StatusInternalServerError, e.NewInternalServerApiError("Unexpected error", err))
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) DeleteProdeById(ctx *gin.Context) {
	prodeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid prode ID"))
		return
	}

	if err := c.prodeService.DeleteProdeById(ctx.Request.Context(), prodeID); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ProdeController) GetProdesByUserId(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
		return
	}

	carreraProdes, sessionProdes, apiErr := c.prodeService.GetProdesByUserId(ctx.Request.Context(), userID)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"carrera_prodes": carreraProdes,
		"session_prodes": sessionProdes,
	})
}

// func (c *ProdeController) GetRaceProdeByUserAndSession(ctx *gin.Context) {
//     userID, err := strconv.Atoi(ctx.Param("user_id"))
//     if err != nil {
//         ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
//         return
//     }

//     sessionID, err := strconv.Atoi(ctx.Param("session_id"))
//     if err != nil {
//         ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
//         return
//     }

//     response, apiErr := c.prodeService.GetRaceProdeByUserAndSession(ctx.Request.Context(), userID, sessionID)
//     if apiErr != nil {
//         ctx.JSON(apiErr.Status(), apiErr)
//         return
//     }

//     ctx.JSON(http.StatusOK, response)
// }

// func (c *ProdeController) GetSessionProdeByUserAndSession(ctx *gin.Context) {
//     userID, err := strconv.Atoi(ctx.Param("user_id"))
//     if err != nil {
//         ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
//         return
//     }

//     sessionID, err := strconv.Atoi(ctx.Param("session_id"))
//     if err != nil {
//         ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
//         return
//     }

//     response, apiErr := c.prodeService.GetSessionProdeByUserAndSession(ctx.Request.Context(), userID, sessionID)
//     if apiErr != nil {
//         ctx.JSON(apiErr.Status(), apiErr)
//         return
//     }

//     ctx.JSON(http.StatusOK, response)
// }

func (c *ProdeController) GetProdeByUserAndSession(ctx *gin.Context) {
    userID, err := strconv.Atoi(ctx.Param("user_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
        return
    }

    sessionID, err := strconv.Atoi(ctx.Param("session_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
        return
    }

    raceProde, sessionProde, apiErr := c.prodeService.GetProdeByUserAndSession(ctx.Request.Context(), userID, sessionID)
    if apiErr != nil {
        // Si es un error diferente a 404, devolvemos el estado del error
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    // Construir la respuesta como un array
    var response []interface{}
    if raceProde != nil {
        response = append(response, *raceProde)
    }
    if sessionProde != nil {
        response = append(response, *sessionProde)
    }
    if len(response) == 0 {
        response = []interface{}{} // Array vacío si no hay datos
    }

    // Devolver siempre un 200 OK con el array (vacío o con datos)
    ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) GetRaceProdesBySession(ctx *gin.Context) {
    sessionID, err := strconv.Atoi(ctx.Param("session_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
        return
    }

    response, apiErr := c.prodeService.GetRaceProdesBySession(ctx.Request.Context(), sessionID)
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) GetSessionProdesBySession(ctx *gin.Context) {
    sessionID, err := strconv.Atoi(ctx.Param("session_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
        return
    }

    response, apiErr := c.prodeService.GetSessionProdeBySession(ctx.Request.Context(), sessionID)
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) UpdateRaceProdeForUserBySessionId(ctx *gin.Context) {
    userID, err := strconv.Atoi(ctx.Param("user_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
        return
    }

    sessionID, err := strconv.Atoi(ctx.Param("session_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
        return
    }

    var request dto.UpdateProdeCarreraDTO
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
        return
    }

    response, apiErr := c.prodeService.UpdateRaceProdeForUserBySessionId(ctx.Request.Context(), userID, sessionID, request)
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) GetUserProdes(ctx *gin.Context) {
    userID, err := strconv.Atoi(ctx.Param("user_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
        return
    }

    carreraProdes, sessionProdes, apiErr := c.prodeService.GetUserProdes(ctx.Request.Context(), userID)
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "carrera_prodes": carreraProdes,
        "session_prodes": sessionProdes,
    })
}

func (c *ProdeController) GetDriverDetails(ctx *gin.Context) {
    driverID, err := strconv.Atoi(ctx.Param("driver_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
        return
    }

    driverDetails, apiErr := c.prodeService.GetDriverDetails(ctx.Request.Context(), driverID)
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, driverDetails)
}

func (c *ProdeController) GetAllDrivers(ctx *gin.Context) {
    drivers, apiErr := c.prodeService.GetAllDrivers(ctx.Request.Context())
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, drivers)
}

func (c *ProdeController) GetTopDriversBySessionId(ctx *gin.Context) {
    sessionID, err := strconv.Atoi(ctx.Param("session_id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid session ID"))
        return
    }

    n, err := strconv.Atoi(ctx.Param("n"))
    if err != nil || n <= 0 {
        ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid number of top drivers"))
        return
    }

    topDrivers, apiErr := c.prodeService.GetTopDriversBySessionId(ctx.Request.Context(), sessionID, n)
    if apiErr != nil {
        ctx.JSON(apiErr.Status(), apiErr)
        return
    }

    ctx.JSON(http.StatusOK, topDrivers)
}