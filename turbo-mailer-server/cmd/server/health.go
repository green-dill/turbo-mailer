package server

import (
	"context"
	"net/http"
	"turbo-mailer-server/internal/query"

	"github.com/labstack/echo/v4"
)

func healthCheck(c echo.Context) error {
	// 检查数据库连接
	sqlDB, err := query.DB.DB()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "error",
			"detail": "Failed to get database connection: " + err.Error(),
		})
	}
	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "error",
			"detail": "Database connection failed: " + err.Error(),
		})
	}

	// 检查Redis连接
	ctx := context.Background()
	if err := query.Redis.Ping(ctx).Err(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "error",
			"detail": "Redis connection failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"detail": "All systems operational",
	})
}
