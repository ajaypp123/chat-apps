package appcontext

import "context"

func DefaultContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "index", "server")
	return ctx
}
