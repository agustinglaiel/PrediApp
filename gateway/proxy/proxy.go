package proxy

// Contiene la implementación del proxy inverso para redirigir las peticiones a los microservicios.

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReverseProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
        target := getTargetURL(c.Request.URL.Path)
		if target == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}

        targetURL, err := url.Parse(target)
		if err != nil {
			log.Printf("Error parsing target URL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host

			// Aquí originalPath representa el path capturado luego del prefijo /users/ (por ejemplo /signup)
			originalPath := c.Param("proxyPath")

			// Asignamos el path tal cual al request hacia el microservicio
			req.URL.Path = originalPath
            
            // Restaurar el body si existe
            if req.Body != nil {
                body, _ := ioutil.ReadAll(req.Body)
                req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
                req.ContentLength = int64(len(body))
            }

			req.Host = targetURL.Host
            log.Printf("Proxying request to: %s %s", req.Method, req.URL.String())
		}

		proxy.ModifyResponse = func(res *http.Response) error {
            body, err := ioutil.ReadAll(res.Body)
            if err != nil {
                return err
            }
            log.Printf("Received response from microservice: %d", res.StatusCode)
            log.Printf("Response body: %s", string(body))
            res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
            return nil
        }

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func getTargetURL(path string) string {
    // Mantenemos un prefijo a nivel de gateway para cada servicio, de este modo:
    // - http://localhost:8080/users/xxx -> http://localhost:8057/xxx (servicio users)
    // - http://localhost:8080/drivers/xxx -> http://localhost:8051/xxx (servicio drivers)
    // y así con los demás servicios
    switch {
    case strings.HasPrefix(path, "/drivers/"):
        return os.Getenv("DRIVERS_SERVICE_URL")
    case strings.HasPrefix(path, "/prodes/"):
        return os.Getenv("PRODES_SERVICE_URL")
    case strings.HasPrefix(path, "/results/"):
        return os.Getenv("RESULTS_SERVICE_URL")
    case strings.HasPrefix(path, "/sessions/"):
        return os.Getenv("SESSIONS_SERVICE_URL")
    case strings.HasPrefix(path, "/users/"):
        return os.Getenv("USERS_SERVICE_URL")
    default:
        return ""
    }
}