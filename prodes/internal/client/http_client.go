package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	dto "prediapp.local/prodes/internal/dto"
	"prediapp.local/prodes/pkg/utils"
)

type CircuitBreakerState string

const (
	Closed       CircuitBreakerState = "CLOSED"
	Open         CircuitBreakerState = "OPEN"
	HalfOpen     CircuitBreakerState = "HALF_OPEN"
	FailLimit                        = 5                // Número máximo de fallos antes de abrir el circuito
	ResetTimeout                     = 10 * time.Second // Tiempo de espera antes de probar si el servicio se ha recuperado
)

type HttpClient struct {
	BaseURL     string
	HTTPClient  *http.Client
	failCount   int
	state       CircuitBreakerState
	mu          sync.Mutex
	lastFailure time.Time
}

// NewHttpClient crea una nueva instancia de HttpClient con una base URL
func NewHttpClient(baseURL string) *HttpClient {
	return &HttpClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		state:      Closed,
	}
}

// checkCircuitState revisa si debe cambiar el estado del circuito
func (c *HttpClient) checkCircuitState() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == Open && time.Since(c.lastFailure) > ResetTimeout {
		c.state = HalfOpen // Intentaremos realizar una solicitud para ver si el servicio se ha recuperado
	}
}

func (c *HttpClient) shouldBlockRequest() bool {
	c.checkCircuitState()
	return c.state == Open
}

func (c *HttpClient) recordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failCount++
	if c.failCount >= FailLimit {
		c.state = Open
		c.lastFailure = time.Now()
	}
}

func (c *HttpClient) resetFailures() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.failCount = 0
	c.state = Closed
}

// Get realiza una solicitud GET a la API de destino con Circuit Breaker
func (c *HttpClient) Get(endpoint string) ([]byte, error) {
	if c.shouldBlockRequest() {
		return nil, fmt.Errorf("circuit breaker is open, blocking request")
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		c.recordFailure()
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.recordFailure()
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	c.resetFailures() // La solicitud fue exitosa, reiniciamos los contadores de fallos

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.recordFailure()
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Post realiza una solicitud POST a la API de destino
func (c *HttpClient) Post(endpoint string, data interface{}) ([]byte, error) {
	if c.shouldBlockRequest() {
		return nil, fmt.Errorf("circuit breaker is open, blocking request")
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request data: %w", err)
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.recordFailure()
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.recordFailure()
		return nil, fmt.Errorf("received non-200/201 response code: %d", resp.StatusCode)
	}

	c.resetFailures()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.recordFailure()
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Put realiza una solicitud PUT a la API de destino
func (c *HttpClient) Put(endpoint string, data interface{}) ([]byte, error) {
	if c.shouldBlockRequest() {
		return nil, fmt.Errorf("circuit breaker is open, blocking request")
	}

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
		c.recordFailure() // Registrar fallo en el Circuit Breaker
		return nil, fmt.Errorf("error making PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.recordFailure() // Registrar fallo en el Circuit Breaker
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	c.resetFailures() // Restablecer contadores de fallo en el Circuit Breaker

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.recordFailure() // Registrar fallo en el Circuit Breaker
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// Delete realiza una solicitud DELETE a la API de destino
func (c *HttpClient) Delete(endpoint string) error {
	if c.shouldBlockRequest() {
		return fmt.Errorf("circuit breaker is open, blocking request")
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.recordFailure() // Registrar fallo en el Circuit Breaker
		return fmt.Errorf("error making DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		c.recordFailure() // Registrar fallo en el Circuit Breaker
		return fmt.Errorf("received non-200/204 response code: %d", resp.StatusCode)
	}

	c.resetFailures() // Restablecer contadores de fallo en el Circuit Breaker

	return nil
}

// GetSessionNameAndType realiza una solicitud GET para obtener el nombre y tipo de una sesión desde el microservicio de sessions
func (c *HttpClient) GetSessionNameAndType(sessionID int) (dto.SessionNameAndTypeDTO, utils.ApiError) {
	endpoint := fmt.Sprintf("/sessions/%d/name-type", sessionID)

	body, err := c.Get(endpoint)
	if err != nil {
		return dto.SessionNameAndTypeDTO{},
			utils.NewInternalServerApiError("Error fetching session name and type", err)
	}

	var dtoResp dto.SessionNameAndTypeDTO
	if err := json.Unmarshal(body, &dtoResp); err != nil {
		return dto.SessionNameAndTypeDTO{},
			utils.NewInternalServerApiError("Error decoding session response", err)
	}
	return dtoResp, nil
}

// GetSessionByID realiza una solicitud GET para obtener los detalles completos de una sesión desde el microservicio de sessions
func (c *HttpClient) GetSessionByID(sessionID int) (dto.SessionDetailsDTO, error) {
	// Montamos la URL completa: ej. "http://sessions:8055/sessions/1"
	endpoint := fmt.Sprintf("/sessions/%d", sessionID)

	// Hacer la solicitud GET utilizando el cliente HTTP
	body, err := c.Get(endpoint)
	if err != nil {
		return dto.SessionDetailsDTO{}, fmt.Errorf("error fetching session by ID: %w", err)
	}

	// Deserializar la respuesta JSON en una estructura Go
	var session dto.SessionDetailsDTO
	if err := json.Unmarshal(body, &session); err != nil {
		return dto.SessionDetailsDTO{}, fmt.Errorf("error decoding session response: %w", err)
	}

	return session, nil
}

func (c *HttpClient) GetUserByID(userID int) (bool, error) {
	endpoint := fmt.Sprintf("/users/%d", userID)
	fmt.Println("Realizando consulta al endpoint: ", endpoint)

	body, err := c.Get(endpoint)
	if err != nil {
		return false, fmt.Errorf("error fetching user by ID: %w", err)
	}

	// Suponemos que la respuesta es un código 200 si el usuario existe, y un 404 si no
	if len(body) == 0 {
		return false, nil
	}

	return true, nil
}

func (c *HttpClient) GetDriverByID(driverID int) (dto.DriverDTO, error) {
	endpoint := fmt.Sprintf("/drivers/%d", driverID)
	body, err := c.Get(endpoint)
	if err != nil {
		return dto.DriverDTO{}, fmt.Errorf("error fetching driver by ID: %w", err)
	}

	var driverDetails dto.DriverDTO
	if err := json.Unmarshal(body, &driverDetails); err != nil {
		return dto.DriverDTO{}, fmt.Errorf("error decoding driver details: %w", err)
	}

	return driverDetails, nil
}

func (c *HttpClient) GetAllDrivers() ([]dto.DriverDTO, error) {
	body, err := c.Get("/drivers")
	if err != nil {
		return nil, fmt.Errorf("error fetching all drivers: %w", err)
	}

	var drivers []dto.DriverDTO
	if err := json.Unmarshal(body, &drivers); err != nil {
		return nil, fmt.Errorf("error decoding drivers list: %w", err)
	}

	return drivers, nil
}

// GetTopDriversBySession realiza una solicitud GET al microservicio de results para obtener los mejores N pilotos de una sesión
func (c *HttpClient) GetTopDriversBySession(sessionID int, n int) ([]dto.TopDriverDTO, error) {
	// Observa que el path puede variar según tu router en Results
	endpoint := fmt.Sprintf("/results/session/%d/top/%d", sessionID, n)
	// Hacer la solicitud GET utilizando el cliente HTTP
	body, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error fetching top drivers: %w", err)
	}

	// Deserializar la respuesta JSON en una lista de TopDriverDTO
	var topDrivers []dto.TopDriverDTO
	if err := json.Unmarshal(body, &topDrivers); err != nil {
		return nil, fmt.Errorf("error decoding top drivers response: %w", err)
	}

	return topDrivers, nil
}
