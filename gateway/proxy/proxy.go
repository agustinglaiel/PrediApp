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
        // log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
        // Manejar las solicitudes preflight directamente
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

        target, proxyPath := getTargetURL(c.Request.URL.Path)
        if target == "" {
            log.Printf("Service not found for path: %s", c.Request.URL.Path)
            c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
            return
        }
        // Si el proxyPath termina en '/', eliminarlo
        if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
            proxyPath = strings.TrimSuffix(proxyPath, "/")
        }
        // log.Printf("Target Base URL: %s, ProxyPath: %s", target, proxyPath)

        targetURL, err := url.Parse(target)
        if err != nil {
            log.Printf("Error parsing target URL: %v", err)

            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            return
        }

        proxy := httputil.NewSingleHostReverseProxy(targetURL)

        proxy.Director = func(req *http.Request) {
            log.Printf("Original Request Path: %s", req.URL.Path)
            log.Printf("Resolved Target URL: %s", targetURL.String())
            log.Printf("Proxy Path: %s", proxyPath)
        
            req.URL.Scheme = targetURL.Scheme
            req.URL.Host = targetURL.Host
            req.URL.Path = proxyPath
        
            if req.Body != nil {
                body, _ := ioutil.ReadAll(req.Body)
                req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
                req.ContentLength = int64(len(body))
            }
        
            req.Host = targetURL.Host
            log.Printf("Forwarding final request to: %s", req.URL.String())
        }

        proxy.ModifyResponse = func(res *http.Response) error {
            if res.StatusCode == http.StatusMovedPermanently || res.StatusCode == http.StatusFound {
                location := res.Header.Get("Location")
                if location != "" {
                    // log.Printf("Redirecting to: %s", location)
                    c.Redirect(res.StatusCode, location) // Redirigir al cliente
                    return nil
                }
            }
            return nil
        }
        

        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func getTargetURL(path string) (string, string) {
    // Quitar el prefijo inicial "/" y dividir el path
    parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

    // Validar que al menos haya un segmento (el microservicio)
    if len(parts) < 1 || parts[0] == "" {
        return "", ""
    }

    service := parts[0] // El primer segmento es el nombre del microservicio
    proxyPath := "/" + strings.Join(parts, "/") // Mantener el path completo

    // Determinar la URL base del microservicio
    switch service {
    case "users":
        return os.Getenv("USERS_SERVICE_URL"), proxyPath
    case "drivers":
        return os.Getenv("DRIVERS_SERVICE_URL"), proxyPath
    case "prodes":
        return os.Getenv("PRODES_SERVICE_URL"), proxyPath
    case "results":
        return os.Getenv("RESULTS_SERVICE_URL"), proxyPath
    case "sessions":
        return os.Getenv("SESSIONS_SERVICE_URL"), proxyPath
    case "groups":
        return os.Getenv("GROUPS_SERVICE_URL"), proxyPath
    default:
        return "", ""
    }
}
