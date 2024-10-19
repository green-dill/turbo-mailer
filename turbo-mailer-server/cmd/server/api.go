package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"turbo-mailer-server/internal/echox"
	"turbo-mailer-server/internal/validator"
	"turbo-mailer-server/version"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/ziflex/lecho/v3"
)

func api() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("500M"))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))
	e.Use(echoprometheus.NewMiddleware("turbo_mailer"))
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = echox.HTTPErrorHandler
	validator.ConfigValidator(e)
	pprof.Register(e)
	e.Logger = lecho.From(log.Logger)

	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s-%s", version.Version, version.CommitID))
	})
	e.GET("/health", healthCheck)

	route(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		addr := viper.GetString("server.address")
		log.Info().Str("address", addr).Msg("Starting server")
		if err := e.Start(addr); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for signal to shutdown
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}
	log.Info().Msg("Server shutdown successfully")
}
