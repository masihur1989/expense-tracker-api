package utils

import "github.com/labstack/echo/v4"

// Response contains properties to be responded
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

// Data returns wrapped success response
func Data(code int, data interface{}, message string, c echo.Context) error {
	props := &Response{
		Code:    code,
		Data:    data,
		Message: message,
		Success: true,
	}
	return c.JSON(code, props)
}

// Error return the wrapped error response
func Error(code int, message string, c echo.Context) error {
	props := &Response{
		Code:    code,
		Data:    nil,
		Message: message,
		Success: false,
	}
	return c.JSON(code, props)
}
