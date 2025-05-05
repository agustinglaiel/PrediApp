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

        // 1. Manejar las preflight requests (OPCIONAL, si ya lo tenías)
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusOK)
            return
        }

        // 2. Normalizar la ruta para eliminar el slash final (p.ej.: "/drivers/" => "/drivers")
        trimmedPath := strings.TrimSuffix(c.Request.URL.Path, "/")
        if trimmedPath == "" {
            trimmedPath = "/" // Evita quedarte sin barra en la raíz
        }
        c.Request.URL.Path = trimmedPath

        // 3. Determinar el microservicio y el path a partir de la ruta
        target, proxyPath := getTargetURL(c.Request.URL.Path)
        if target == "" {
            log.Printf("Service not found for path: %s", c.Request.URL.Path)
            c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
            return
        }

        // 4. Eliminar barra final también en proxyPath (por si acaso)
        if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
            proxyPath = strings.TrimSuffix(proxyPath, "/")
        }

        // 5. Parsear la URL base del microservicio (ej. "http://localhost:8051")
        targetURL, err := url.Parse(target)
        if err != nil {
            log.Printf("Error parsing target URL: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            return
        }

        // 6. Crear el Reverse Proxy de Go
        proxy := httputil.NewSingleHostReverseProxy(targetURL)

        // Director: manipula la request ANTES de mandarla al microservicio
        proxy.Director = func(req *http.Request) {
            log.Printf("Original Request Path: %s", req.URL.Path)
            log.Printf("Resolved Target URL: %s", targetURL.String())
            log.Printf("Proxy Path: %s", proxyPath)

            // Ajustamos la request al target final
            req.URL.Scheme = targetURL.Scheme
            req.URL.Host = targetURL.Host
            req.URL.Path = proxyPath

            // En caso de que venga body, se vuelve a leer
            if req.Body != nil {
                body, _ := ioutil.ReadAll(req.Body)
                req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
                req.ContentLength = int64(len(body))
            }

            // Ajustamos el Host
            req.Host = targetURL.Host
            log.Printf("Forwarding final request to: %s", req.URL.String())
        }

        // ModifyResponse (OPCIONAL): Para manejar redirecciones o modificar la respuesta
        proxy.ModifyResponse = func(res *http.Response) error {
            // Si el microservicio manda un 301/302/307...
            if res.StatusCode == http.StatusMovedPermanently || res.StatusCode == http.StatusFound {
                location := res.Header.Get("Location")
                if location != "" {
                    c.Redirect(res.StatusCode, location)
                    return nil
                }
            }
            return nil
        }

        // 7. Finalmente, servir la request al microservicio
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
