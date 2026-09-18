package bootstrap

import (
	"context"
	"gobnb/internal/configuration"
)

func Initialize() context.Context {
	ctx := context.Background()
	ctx.Deadline()
	configuration.Load()
	return ctx
}
