package mailer

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Domains(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
