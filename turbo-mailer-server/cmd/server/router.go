package server

import (
	"net/http"
	"turbo-mailer-server/internal/api/admin/auth"
	"turbo-mailer-server/internal/api/admin/dashboard"
	"turbo-mailer-server/internal/api/admin/pool"
	"turbo-mailer-server/internal/api/admin/task"

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

	g.GET("/task", task.List)
	g.GET("/task/:id", task.Get)
	g.GET("/tasks/content-template", task.DownloadContentTemplate)
	g.GET("/tasks/receivers-template", task.DownloadReceiversTemplate)
	g.POST("/task", task.Store)
	g.POST("/task/:id", task.Update)
	g.DELETE("/task/:id", task.Delete)
	g.POST("/task/:id/test", task.Test)
	g.POST("/task/:id/start-immediately", task.StartImmediately)

	g.GET("/mail-server-dns", todo)
}

func todo(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "todo")
}
