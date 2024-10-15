package dispatch

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func Boot(ctx context.Context) {
	go run(ctx)
}

func run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("panic: %v", r)
		}
		time.Sleep(time.Second)
		run(ctx)
	}()

	dispatcher, err := newDispatcher()
	if err != nil {
		log.Error().Err(err).Msg("failed to create dispatcher")
		return
	}
	dispatcher.Run(ctx)
}
