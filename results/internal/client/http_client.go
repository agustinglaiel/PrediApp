package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	dto "results/internal/dto"
	"time"
)

type HttpClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewHttpClient crea una nueva instancia de HttpClient con una base URL
func NewHttpClient(baseURL string) *HttpClient {
	return &HttpClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Get realiza una solicitud GET a la API de destino
func (c *HttpClient) Get(endpoint string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Post realiza una solicitud POST a la API de destino
func (c *HttpClient) Post(endpoint string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request data: %w", err)
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("received non-200/201 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Put realiza una solicitud PUT a la API de destino
func (c *HttpClient) Put(endpoint string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Delete realiza una solicitud DELETE a la API de destino
func (c *HttpClient) Delete(endpoint string) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received non-200/204 response code: %d", resp.StatusCode)
	}

	return nil
}

// GetPositions obtiene las posiciones de los pilotos desde la API externa para una sesión específica
func (c *HttpClient) GetPositions(sessionKey int) ([]dto.Position, error) {
    // Construir la URL completa para la solicitud de posiciones
    endpoint := fmt.Sprintf("https://api.openf1.org/v1/position?session_key=%d", sessionKey)
	println("Endpoint: ", endpoint)
    fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint) // BaseURL ya tiene el esquema y dominio

    body, err := c.Get(fullURL)
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
    fullURL := fmt.Sprintf("%s%s", c.BaseURL, endpoint) // BaseURL ya tiene el esquema y dominio

    body, err := c.Get(fullURL)
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
	// Usar la URL correcta del microservicio de sessions
	endpoint := fmt.Sprintf("http://localhost:8056/sessions/%d/get-session-key", sessionID)
	fmt.Println("Endpoint: ", endpoint)

	// Hacer la solicitud GET utilizando el cliente HTTP
	body, err := c.Get(endpoint)
	if err != nil {
		return 0, fmt.Errorf("error fetching sessionKey: %w", err)
	}

	// Deserializar la respuesta para obtener el sessionKey
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
	endpoint := fmt.Sprintf("http://localhost:8051/drivers/number/%d", driverNumber)

	// Hacer la solicitud GET utilizando el cliente HTTP
	body, err := c.Get(endpoint)
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

	// Imprimir el driver_number y el driver_id que se está manejando
	fmt.Printf("Obtenido driver_number: %d, driver_id: %d desde el microservicio drivers\n", driverNumber, driver.ID)

	return driver, nil
}

// Función para obtener la información de una sesión completa utilizando sessionId
func (c *HttpClient) GetSessionByID(sessionID int) (dto.ResponseSessionDTO, error) {
    // Usar la URL correcta del microservicio de sessions
    endpoint := fmt.Sprintf("http://localhost:8056/sessions/%d", sessionID)

    // Hacer la solicitud GET utilizando el cliente HTTP
    body, err := c.Get(endpoint)
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
