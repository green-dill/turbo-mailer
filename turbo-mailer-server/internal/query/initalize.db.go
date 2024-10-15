package query

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB    *gorm.DB
	Redis *redisWarp

	initializeOnce sync.Once
)

func Initialize(context.Context) error {
	initializeOnce.Do(func() {
		initializeDb()
		initializeRedis()
	})
	return nil
}

type DbConfig struct {
	Uri                          string `yaml:"uri"`
	MaxConn                      int    `yaml:"maxConn"`
	MaxIdleConn                  int    `yaml:"maxIdleConn"`
	ConnMaxLifetimeInMillisecond int    `yaml:"connMaxLifetimeInMillisecond"`
	QueryFields                  bool   `yaml:"queryFields"`
	CreateBatchSize              int    `yaml:"createBatchSize"`
}

func initializeDb() {
	dbMap := make(map[string]*DbConfig)

	if err := viper.UnmarshalKey("database", &dbMap); err != nil {
		log.Fatal().Err(err).Msg("parse database config")
	}

	// set default value
	for _, cfg := range dbMap {
		if cfg.MaxConn == 0 {
			cfg.MaxConn = 10
		}
		if cfg.MaxIdleConn == 0 {
			cfg.MaxIdleConn = 10
		}
		if cfg.ConnMaxLifetimeInMillisecond == 0 {
			cfg.ConnMaxLifetimeInMillisecond = 60000
		}
		if cfg.CreateBatchSize == 0 {
			cfg.CreateBatchSize = 1000
		}
	}

	if db, err := buildOrm(dbMap["llm"]); err != nil {
		log.Fatal().Err(err).Msg("init db")
	} else {
		DB = db
	}

	SetDefault(DB)
}

func buildOrm(cfg *DbConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.Uri), &gorm.Config{
		QueryFields:                              cfg.QueryFields,
		CreateBatchSize:                          cfg.CreateBatchSize,
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: &gormLogger{},
	})
}

func initializeRedis() {
	var cfgMap map[string]*redisConfig

	if err := viper.UnmarshalKey("redis", &cfgMap); err != nil {
		log.Fatal().Err(err).Msg("parse redis config")
	}
	if c, ok := cfgMap["redis"]; ok {
		Redis = newRedis(c)
	}
}

type gormLogger struct{}

// Error implements logger.Interface.
func (l *gormLogger) Error(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Error().Msgf(msg, args...)
}

// Info implements logger.Interface.
func (l *gormLogger) Info(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Info().Msgf(msg, args...)
}

// LogMode implements logger.Interface.
func (l *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

// Trace implements logger.Interface.
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	end := time.Now()
	escaped := end.Sub(begin)
	log.Ctx(ctx).Debug().Msgf("sql: %s, rowsAffected: %d, escaped: %s, err: %v", sql, rowsAffected, escaped, err)
}

// Warn implements logger.Interface.
func (l *gormLogger) Warn(ctx context.Context, msg string, args ...any) {
	log.Ctx(ctx).Warn().Msgf(msg, args...)
}
