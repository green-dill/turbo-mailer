package initialize

import (
	"context"
	"reflect"
	"runtime"
	"sync"
	"turbo-mailer-server/internal/query"

	"github.com/rs/zerolog/log"
)

type initializeFn func(ctx context.Context) error

var (
	once sync.Once

	initializeFns = []initializeFn{
		query.Initialize,
	}
)

func Do(ctx context.Context) {
	once.Do(func() {
		for _, fn := range initializeFns {
			if err := fn(ctx); err != nil {
				fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
				log.Fatal().Err(err).Str("function", fnName).Msg("failed to initialize")
			}
		}
	})
}
