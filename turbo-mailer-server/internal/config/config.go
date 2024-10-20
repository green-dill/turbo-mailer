package config

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/spf13/viper"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.TimeFieldFormat = time.DateTime
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: time.DateTime})
	// log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = log.With().Caller().Logger()

	zerolog.DefaultContextLogger = &log.Logger

	slog.SetDefault(slog.New(slogzerolog.Option{Level: slog.LevelError}.NewZerologHandler()))
}

func Initialize(configPath ...string) {
	log.Debug().Msg("initialize config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	for _, path := range configPath {
		viper.AddConfigPath(path)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("read config file error")
	}

	viper.SetDefault("server.address", "0.0.0.0:5000")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("jwt.secret", "secret")

	if lvl, err := zerolog.ParseLevel(viper.GetString("log.level")); err == nil {
		log.Logger = log.Logger.Level(lvl)
	} else {
		log.Error().Err(err).Msg("parse log level error")
	}
}
