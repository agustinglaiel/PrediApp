package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	dto "results/internal/dto"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(secret string, userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Expira en 1 hora
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type HttpClient struct {
	BaseURL    string
	HTTPClient *http.Client
	JWTToken   string
}

// NewHttpClient crea una nueva instancia de HttpClient con una base URL
func NewHttpClient(baseURL string) *HttpClient {
    if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
        panic("BaseURL must start with http:// or https://")
    }
    return &HttpClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: time.Second * 30,
        },
    }
}


func (c *HttpClient) SetJWTToken(token string) {
	c.JWTToken = token
}

// Get realiza una solicitud GET a la API de destino
func (c *HttpClient) Get(endpoint string) ([]byte, error) {
    // Verificar si el endpoint ya es una URL completa
    url := endpoint
    if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
        url = fmt.Sprintf("%s%s", c.BaseURL, endpoint)
    }

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating GET request: %w", err)
    }

    // Agregar el token JWT si está configurado
    if c.JWTToken != "" {
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWTToken))
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error making GET request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("received non-200 response code: %d, body: %s", resp.StatusCode, body)
    }

    return ioutil.ReadAll(resp.Body)
}



// Post realiza una solicitud POST a la API de destino
func (c *HttpClient) Post(endpoint string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Agregar el token JWT si está configurado
	if c.JWTToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWTToken))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200/201 response code: %d, body: %s", resp.StatusCode, body)
	}

	return ioutil.ReadAll(resp.Body)
}


// Put realiza una solicitud PUT a la API de destino
func (c *HttpClient) Put(endpoint string, data interface{}) ([]byte, error) {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	// Serializar los datos a JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request data: %w", err)
	}

	// Crear la solicitud PUT
	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Agregar el token JWT al encabezado
	if c.JWTToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWTToken))
	}

	// Hacer la solicitud PUT
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200/204 response code: %d, body: %s", resp.StatusCode, body)
	}

	// Leer la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}


// Delete realiza una solicitud DELETE a la API de destino
func (c *HttpClient) Delete(endpoint string) error {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	// Crear la solicitud DELETE
	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %w", err)
	}

	// Agregar el token JWT al encabezado
	if c.JWTToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWTToken))
	}

	// Hacer la solicitud DELETE
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-200/204 response code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}


// GetPositions obtiene las posiciones de los pilotos desde la API externa para una sesión específica
func (c *HttpClient) GetPositions(sessionKey int) ([]dto.Position, error) {
    // Construir la URL completa para la solicitud de posiciones
    endpoint := fmt.Sprintf("https://api.openf1.org/v1/position?session_key=%d", sessionKey)
	println("Endpoint: ", endpoint)

    body, err := c.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("error fetching positions: %w", err)
    }

    var positions []dto.Position
    if err := json.Unmarshal(body, &positions); err != nil {
        return nil, fmt.Errorf("error decoding positions response: %w", err)
    }

    return positions, nil
}

// GetLaps obtiene las vueltas rápidas de un piloto específico desde la API externa
func (c *HttpClient) GetLaps(sessionKey int, driverNumber int) ([]dto.Lap, error) {
    // Construir la URL completa para la solicitud de laps
    endpoint := fmt.Sprintf("https://api.openf1.org/v1/laps?session_key=%d&driver_number=%d", sessionKey, driverNumber)
	println("Endpoint: ", endpoint)

    body, err := c.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("error fetching laps: %w", err)
    }

    var laps []dto.Lap
    if err := json.Unmarshal(body, &laps); err != nil {
        return nil, fmt.Errorf("error decoding laps response: %w", err)
    }

    // Filtrar las vueltas con lap_duration nulo o igual a 0
    validLaps := make([]dto.Lap, 0)
    for _, lap := range laps {
        if lap.LapDuration > 0 {
            validLaps = append(validLaps, lap)
        }
    }

    return validLaps, nil
}

// Función para obtener sessionKey utilizando sessionId
func (c *HttpClient) GetSessionKeyBySessionID(sessionID int) (int, error) {
	endpoint := c.buildURL(fmt.Sprintf("/sessions/%d/get-session-key", sessionID))

	body, err := c.GetWithAuth(endpoint)
	if err != nil {
		return 0, fmt.Errorf("error fetching sessionKey: %w", err)
	}

	var response struct {
		SessionKey int `json:"session_key"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("error decoding sessionKey response: %w", err)
	}

	return response.SessionKey, nil
}


// GetDriverByNumber obtiene la información de un piloto basado en su driver_number
func (c *HttpClient) GetDriverByNumber(driverNumber int) (dto.ResponseDriverDTO, error) {
	// Definir el endpoint para obtener la información del piloto desde el driver_number
	endpoint := c.buildURL(fmt.Sprintf("/drivers/number/%d", driverNumber))

	// Hacer la solicitud GET utilizando el cliente HTTP con autenticación
	body, err := c.GetWithAuth(endpoint) // Usar la función que maneja el token JWT
	if err != nil {
		return dto.ResponseDriverDTO{}, fmt.Errorf("error fetching driver info: %w", err)
	}

	// Declarar la variable para deserializar la respuesta
	var driver dto.ResponseDriverDTO

	// Deserializar la respuesta para obtener la información del piloto
	if err := json.Unmarshal(body, &driver); err != nil {
		return dto.ResponseDriverDTO{}, fmt.Errorf("error decoding driver response: %w", err)
	}

	// Verificar si el driver ID es válido
	if driver.ID == 0 {
		return dto.ResponseDriverDTO{}, fmt.Errorf("driver with number %d not found or has invalid ID", driverNumber)
	}

	return driver, nil
}

// Función para obtener la información de una sesión completa utilizando sessionId
func (c *HttpClient) GetSessionByID(sessionID int) (dto.ResponseSessionDTO, error) {
    // Usar la URL correcta del microservicio de sessions
    endpoint := c.buildURL(fmt.Sprintf("/sessions/%d", sessionID))

    // Hacer la solicitud GET utilizando el cliente HTTP con autenticación
    body, err := c.GetWithAuth(endpoint) // Usar la función que maneja el token JWT
    if err != nil {
        return dto.ResponseSessionDTO{}, fmt.Errorf("error fetching session by ID: %w", err)
    }

    // Deserializar la respuesta para obtener la sesión
    var session dto.ResponseSessionDTO
    if err := json.Unmarshal(body, &session); err != nil {
        return dto.ResponseSessionDTO{}, fmt.Errorf("error decoding session response: %w", err)
    }

    return session, nil
}

func (c *HttpClient) GetWithAuth(endpoint string) ([]byte, error) {
	// Crear la solicitud GET directamente con el endpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Obtener la clave secreta JWT del entorno
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is missing")
	}

	// Generar el token JWT dinámico usando GenerateJWT
	token, err := GenerateJWT(secretKey, 1, "service") // Usuario y rol fijo para este cliente
	if err != nil {
		return nil, fmt.Errorf("error generating JWT: %w", err)
	}

	// Configurar el encabezado Authorization con el token JWT
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Hacer la solicitud
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Validar el código de respuesta
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 response code: %d, body: %s", resp.StatusCode, body)
	}

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}


func (c *HttpClient) buildURL(endpoint string) string {
	// Si el endpoint ya es una URL completa, devolverlo tal cual
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		return endpoint
	}
	// Si es relativo, combinarlo con BaseURL
	return fmt.Sprintf("%s%s", strings.TrimRight(c.BaseURL, "/"), endpoint)
}