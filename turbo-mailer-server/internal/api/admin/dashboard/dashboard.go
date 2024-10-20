package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Stats godoc
//
//	@Summary		Get dashboard statistics
//	@Description	Retrieve statistics for the admin dashboard
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"ok"
//	@Router			/api/v1/dashboard/stats [get]
func Stats(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
