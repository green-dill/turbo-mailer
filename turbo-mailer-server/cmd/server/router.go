package server

import (
	"net/http"
	"turbo-mailer-server/internal/api/admin/auth"
	"turbo-mailer-server/internal/api/admin/dashboard"
	"turbo-mailer-server/internal/api/admin/pool"

	"github.com/labstack/echo/v4"
)

func route(e *echo.Echo) {
	g := e.Group("/api/v1")

	g.POST("/auth/login", auth.Login)
	g.POST("/auth/logout", auth.Logout)
	g.POST("/auth/change-password", auth.ChangePassword)
	g.GET("/auth/profile", auth.Profile)

	g.GET("/dashboard/stats", dashboard.Stats)

	g.GET("/pool", pool.List)
	g.GET("/pool/:id", pool.Get)
	g.POST("/pool", pool.Store)
	g.POST("/pool/:id", pool.Update)
	g.DELETE("/pool/:id", pool.Delete)
	g.GET("/pool-senders/:id", pool.SenderList)
	g.POST("/pool-senders/:id", pool.SenderStore)
	g.DELETE("/pool-senders/:sid", pool.SenderDelete)

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
