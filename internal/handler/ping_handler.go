package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/masihur1989/expense-tracker-api/internal/utils"
)

// Ping godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response
// @Router / [get]
func Ping(c echo.Context) error {
	return utils.Data(http.StatusOK, nil, "Server is Up", c)
}
