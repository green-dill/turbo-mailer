package worker

import (
	"context"
)

func Boot(ctx context.Context) {
	go func() {
		run(ctx)
	}()
}
