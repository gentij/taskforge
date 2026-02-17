package cli

import (
	"context"

	"github.com/gentij/taskforge/apps/cli/internal/api"
	"github.com/gentij/taskforge/apps/cli/internal/config"
)

type Context struct {
	Config     config.Config
	Client     *api.Client
	OutputJSON bool
	Quiet      bool
}

type ctxKey struct{}

func WithContext(ctx context.Context, value *Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, value)
}

func GetContext(ctx context.Context) *Context {
	value := ctx.Value(ctxKey{})
	if value == nil {
		return nil
	}

	if typed, ok := value.(*Context); ok {
		return typed
	}

	return nil
}
