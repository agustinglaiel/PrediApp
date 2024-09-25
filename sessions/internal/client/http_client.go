package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	dto "sessions/internal/dto"
	"sessions/pkg/utils"
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

// GetRaceControlData realiza una solicitud GET para obtener los datos de control de carrera (VSC/SC) de la API externa
func (c *HttpClient) GetRaceControlData(sessionKey *int) ([]dto.RaceControlEvent, error) {
    // Definir el endpoint para la solicitud GET
    endpoint := fmt.Sprintf("race_control?session_key=%d", *sessionKey)

    // Hacer la solicitud GET utilizando el cliente HTTP
    body, err := c.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("error fetching race control data: %w", err)
    }

    // Deserializar la respuesta JSON en una estructura Go
    var raceControlEvents []dto.RaceControlEvent
    if err := json.Unmarshal(body, &raceControlEvents); err != nil {
        return nil, fmt.Errorf("error decoding race control response: %w", err)
    }

    return raceControlEvents, nil
}

/*
// GetLapsData obtiene las vueltas de los pilotos para una sesión específica
func (c *HttpClient) GetLapsData(sessionKey int) ([]dto.LapData, e.ApiError) {
    // Definir el endpoint para la solicitud GET de las vueltas
    endpoint := fmt.Sprintf("/laps?session_key=%d", sessionKey)

    // Hacer la solicitud GET utilizando el cliente HTTP
    body, err := c.Get(endpoint)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error al obtener los datos de las vueltas", err)
    }

    // Deserializar la respuesta JSON en una estructura Go
    var lapsData []dto.LapData
    if err := json.Unmarshal(body, &lapsData); err != nil {
        return nil, e.NewInternalServerApiError("Error al decodificar la respuesta de las vueltas", err)
    }

    // Retornar los datos de las vueltas y nil en caso de éxito
    return lapsData, nil
}*/

//Esta función la usamos para obtener el session_key de una session para luego poder hacer el update de sc y vsc
// GetSessionKey obtiene el session_key basado en location, session_name, session_type, y year
// GetSessionData obtiene session_key, date_start y date_end basado en location, session_name, session_type y year
func (c *HttpClient) GetSessionData(location string, sessionName string, sessionType string, year int) (*dto.SessionKeyResponseDTO, utils.ApiError) {
    // Definir el endpoint con los parámetros
    endpoint := fmt.Sprintf("/sessions?location=%s&session_name=%s&session_type=%s&year=%d", location, sessionName, sessionType, year)
    
    // Print del endpoint que estamos consultando
    fmt.Printf("Consultando endpoint: %s\n", endpoint)
    
    // Hacer la solicitud GET
    body, err := c.Get(endpoint)
    if err != nil {
        fmt.Printf("Error en la solicitud GET: %v\n", err)
        return nil, utils.NewInternalServerApiError("Error fetching session data", err)
    }
    
    // Print del body de la respuesta (para ver qué devuelve la API)
    fmt.Printf("Respuesta del body: %s\n", string(body))

    // Deserializar la respuesta JSON como un array de SessionKeyResponseDTO
    var response []dto.SessionKeyResponseDTO
    if err := json.Unmarshal(body, &response); err != nil {
        fmt.Printf("Error deserializando la respuesta JSON: %v\n", err)
        return nil, utils.NewInternalServerApiError("Error decoding session data response", err)
    }

    // Print de la respuesta deserializada para ver qué contiene el objeto `response`
    fmt.Printf("Respuesta deserializada: %+v\n", response)

    // Verificar si el array tiene al menos un resultado y si el session_key está presente
    if len(response) == 0 || response[0].SessionKey == nil {
        fmt.Println("No se encontraron datos de session_key para los parámetros proporcionados.")
        return nil, utils.NewNotFoundApiError("Session data not found for the given parameters")
    }

    // Print de los datos que estamos retornando
    fmt.Printf("Datos retornados: SessionKey: %d, DateStart: %v, DateEnd: %v\n", *response[0].SessionKey, response[0].DateStart, response[0].DateEnd)

    // Retornar el objeto con session_key, date_start y date_end del primer resultado
    return &response[0], nil
}

// GetFastestLapBySessionID obtiene el piloto con la vuelta más rápida de una sesión específica
func (c *HttpClient) GetFastestLapBySessionID(sessionID uint) (*dto.FastestLapDTO, error) {
    endpoint := fmt.Sprintf("http://localhost:8071/results/session/%d/fastest-lap", sessionID)
    fmt.Println("Requesting:", endpoint)  // Log de la solicitud

    body, err := c.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("error fetching fastest lap: %w", err)
    }

    if len(body) == 0 {
        return nil, fmt.Errorf("error fetching fastest lap: empty response")
    }

    // Deserializar la respuesta
    var fastestLap dto.FastestLapDTO
    if err := json.Unmarshal(body, &fastestLap); err != nil {
        return nil, fmt.Errorf("error decoding fastest lap response: %w", err)
    }

    return &fastestLap, nil
}