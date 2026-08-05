package bootstrap

import (
	"GoBNB/internal/configuration"
	"context"
)

func Initialize() context.Context {
	ctx := context.Background()
	ctx.Deadline()
	configuration.Load()
	return ctx
}
