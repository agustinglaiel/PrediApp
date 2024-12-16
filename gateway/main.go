package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.Default()

    // Mantenemos un prefijo a nivel de gateway para cada servicio, de este modo:
    // - http://localhost:8080/users/xxx -> http://localhost:8057/xxx (servicio users)
    // - http://localhost:8080/drivers/xxx -> http://localhost:8051/xxx (servicio drivers)
    // y así con los demás servicios
    router.Any("/drivers/*proxyPath", reverseProxy("http://localhost:8051"))
    router.Any("/prodes/*proxyPath", reverseProxy("http://localhost:8054"))
    router.Any("/results/*proxyPath", reverseProxy("http://localhost:8055"))
    router.Any("/sessions/*proxyPath", reverseProxy("http://localhost:8056"))
    router.Any("/users/*proxyPath", reverseProxy("http://localhost:8057"))

    log.Printf("Gateway escuchando en el puerto %s", port)
    router.Run(":" + port)
}

func reverseProxy(target string) gin.HandlerFunc {
    return func(c *gin.Context) {
        targetURL, _ := url.Parse(target)
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
