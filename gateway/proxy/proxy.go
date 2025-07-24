package proxy

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
		// 0) Quitar el prefijo "/api" para que el path interno empiece con "/users", "/sessions", etc.
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api")

		// 1) Manejar preflight CORS
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// 2) Normalizar la ruta eliminando slash final
		trimmed := strings.TrimSuffix(c.Request.URL.Path, "/")
		if trimmed == "" {
			trimmed = "/"
		}
		c.Request.URL.Path = trimmed

		// 3) Determinar el microservicio y la ruta proxy
		target, proxyPath := getTargetURL(c.Request.URL.Path)
		if target == "" {
			log.Printf("Service not found for path: %s", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}

		// 4) Volver a normalizar proxyPath
		if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
			proxyPath = strings.TrimSuffix(proxyPath, "/")
		}

		// 5) Parsear la URL del target
		targetURL, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// 6) Crear y configurar el reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = proxyPath

			// Si hay body, lo volvemos a leer
			if req.Body != nil {
				b, _ := ioutil.ReadAll(req.Body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
				req.ContentLength = int64(len(b))
			}

			req.Host = targetURL.Host
		}

		proxy.ModifyResponse = func(res *http.Response) error {
			if res.StatusCode == http.StatusMovedPermanently || res.StatusCode == http.StatusFound {
				if loc := res.Header.Get("Location"); loc != "" {
					c.Redirect(res.StatusCode, loc)
				}
			}
			return nil
		}

		// 7) Servir la petición al microservicio
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// getTargetURL se mantiene igual que antes (o con tu lógica actual)
func getTargetURL(path string) (string, string) {
	// Quitar prefijo inicial "/" y dividir
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(parts) < 1 || parts[0] == "" {
		return "", ""
	}

	service := parts[0]
	proxyPath := "/" + strings.Join(parts, "/")

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
	case "posts":
		return os.Getenv("POSTS_SERVICE_URL"), proxyPath
	default:
		return "", ""
	}
}
