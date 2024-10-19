package task

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func Get(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func Store(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func Delete(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
