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
			Timeout: time.Second * 10,
		},
	}
}

// Get realiza una solicitud GET a la API de destino
func (c *HttpClient) Get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	resp, err := c.HTTPClient.Get(url)
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
    endpoint := fmt.Sprintf("position?session_key=%d", sessionKey)
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
    endpoint := fmt.Sprintf("laps?session_key=%d&driver_number=%d", sessionKey, driverNumber)
    body, err := c.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("error fetching laps: %w", err)
    }

    var laps []dto.Lap
    if err := json.Unmarshal(body, &laps); err != nil {
        return nil, fmt.Errorf("error decoding laps response: %w", err)
    }

    return laps, nil
}

// Función para obtener sessionKey utilizando sessionId
func (c *HttpClient) GetSessionKeyBySessionID(sessionID uint) (int, error) {
	// Definir el endpoint para obtener la sessionKey desde el sessionID
	endpoint := fmt.Sprintf("http://localhost:8060/sessions/%d/get-session-key", sessionID)

	// Hacer la solicitud GET utilizando el cliente HTTP
	body, err := c.Get(endpoint)
	if err != nil {
		return 0, fmt.Errorf("error fetching sessionKey: %w", err)
	}

	// Deserializar la respuesta para obtener el sessionKey (esto depende de la estructura de la respuesta)
	var response struct {
		SessionKey int `json:"session_key"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("error decoding sessionKey response: %w", err)
	}

	return response.SessionKey, nil
}
