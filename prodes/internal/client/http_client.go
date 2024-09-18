package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	dto "prodes/internal/dto"
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

// GetSessionNameAndType realiza una solicitud GET para obtener el nombre y tipo de una sesión desde el microservicio de sessions
func (c *HttpClient) GetSessionNameAndType(eventID int) (dto.SessionNameAndTypeDTO, error) {
    // Definir el endpoint para la solicitud GET
    endpoint := fmt.Sprintf("/sessions/%d/name-type", eventID)

    // Hacer la solicitud GET utilizando el cliente HTTP
    body, err := c.Get(endpoint)
    if err != nil {
        return dto.SessionNameAndTypeDTO{}, fmt.Errorf("error fetching session name and type: %w", err)
    }

    // Deserializar la respuesta JSON en una estructura Go
    var sessionNameAndType dto.SessionNameAndTypeDTO
    if err := json.Unmarshal(body, &sessionNameAndType); err != nil {
        return dto.SessionNameAndTypeDTO{}, fmt.Errorf("error decoding session name and type response: %w", err)
    }

    return sessionNameAndType, nil
}