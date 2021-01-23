package middleware

import (
	"log"
	"net/http/httputil"
	"strings"

	"github.com/labstack/echo/v4"
)

// RequestHeaders log the all the request header
func RequestHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			httpRequest, _ := httputil.DumpRequest(c.Request(), false)
			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					// mask the auth token
					headers[idx] = current[0] + ": *"
				}
			}
			log.Printf("New Request: \n%s", strings.Join(headers, "\r\n"))
			return next(c)
		}
	}
}
