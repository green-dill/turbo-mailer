package echox

import (
	"errors"
	"fmt"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if he, ok := err.(*echo.HTTPError); ok {
		switch err {
		case middleware.ErrJWTMissing:
			_ = c.JSON(http.StatusUnauthorized, echo.Map{"message": "login required"})
			return
		default:
			status := he.Code
			message := fmt.Sprintf("%v", he.Message)
			if he.Internal != nil {
				message = fmt.Sprintf("%v, cause: %v", he.Message, he.Internal)
			}
			_ = c.JSON(status, echo.Map{"message": message})
			return
		}
	} else if ve, ok := err.(validator.ValidationErrors); ok {
		_ = c.JSON(http.StatusBadRequest, echo.Map{"message": ve.Error()})
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = c.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
		return
	} else {
		_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":     err.Error(),
			"requestId": c.Get(echo.HeaderXRequestID),
		})
	}
}
