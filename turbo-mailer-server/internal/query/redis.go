package query

import (
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisConfig struct {
	Name         string `yaml:"name"`
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	ReadTimeout  int    `yaml:"readTimeout"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"poolSize"`
	MaxRetries   int    `yaml:"maxRetries"`
	MinIdleConns int    `yaml:"minIdleConns"`
	MaxConnAge   int    `yaml:"maxConnAge"`
	Prefix       string `yaml:"prefix"`
}

type redisWarp struct {
	*redis.Client
	cfg *redisConfig
}

func newRedis(c *redisConfig) *redisWarp {
	client := redis.NewClient(&redis.Options{
		Addr:            c.Addr,
		Password:        c.Password,
		DB:              c.DB,
		ReadTimeout:     time.Second * time.Duration(c.ReadTimeout),
		MaxRetries:      c.MaxRetries,
		MinIdleConns:    c.MinIdleConns,
		ConnMaxLifetime: time.Second * time.Duration(c.MaxConnAge),
		PoolSize:        c.PoolSize,
	})
	return &redisWarp{Client: client, cfg: c}
}

func (r *redisWarp) Key(key ...string) string {
	return strings.Join(append([]string{r.cfg.Prefix}, key...), ":")
}
