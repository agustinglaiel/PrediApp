package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// UsersClient maneja las solicitudes al microservicio de users
type UsersClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewUsersClient crea una nueva instancia del cliente
func NewUsersClient() *UsersClient {
	baseURL := os.Getenv("USERS_SERVICE_URL") // URL del microservicio de users
	if baseURL == "" {
		panic("USERS_SERVICE_URL environment variable is missing")
	}

	return &UsersClient{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// GetUserScore obtiene el puntaje de un usuario desde el microservicio de users
func (c *UsersClient) GetUserScore(userID int) (int, error) {
	url := fmt.Sprintf("%s/users/%d", c.BaseURL, userID) // Endpoint actualizado

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating GET request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("received non-200 response code: %d, body: %s", resp.StatusCode, body)
	}

	// Definir la estructura completa para decodificar la respuesta
	var response struct {
		Score int `json:"score"`
	}

	// Decodificar la respuesta
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding user response: %w", err)
	}

	return response.Score, nil
}

