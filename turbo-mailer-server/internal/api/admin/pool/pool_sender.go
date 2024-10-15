package pool

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SenderList(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func SenderStore(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}

func SenderDelete(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
