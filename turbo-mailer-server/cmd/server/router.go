package server

import (
	"net/http"
	"turbo-mailer-server/internal/api/admin/auth"

	"github.com/labstack/echo/v4"
)

func route(e *echo.Echo) {
	g := e.Group("/api/v1")

	g.POST("/auth/login", auth.Login)
	g.POST("/auth/logout", auth.Logout)
	g.POST("/auth/change-password", auth.ChangePassword)

	g.GET("/dashboard/stats", todo)

	g.GET("/pool", todo)
	g.GET("/pool/:id", todo)
	g.POST("/pool", todo)
	g.PUT("/pool/:id", todo)
	g.DELETE("/pool/:id", todo)
	g.POST("/pool/:id", todo)
	g.PUT("/pool/:id/:sid", todo)
	g.DELETE("/pool/:id/:sid", todo)

	g.GET("/task", todo)
	g.GET("/task/:id", todo)
	g.POST("/task", todo)
	g.PUT("/task/:id", todo)
	g.DELETE("/task/:id", todo)
	g.POST("/task/:id/test", todo)
	g.POST("/task/:id/start-immediately", todo)

	g.GET("/mail-server-dns", todo)
}

func todo(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "todo")
}
