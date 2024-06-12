package main

import (
	"context"
	"os"

	"github.com/abibby/remote-input/server/app"
	"github.com/abibby/salusa/di"
)

func main() {
	ctx := di.ContextWithDependencyProvider(
		context.Background(),
		di.NewDependencyProvider(),
	)

	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		app.Kernel.Logger(ctx).Error("error bootstrapping", "error", err)
		os.Exit(1)
	}

	err = app.Kernel.Run(ctx)
	if err != nil {
		app.Kernel.Logger(ctx).Error("error running", "error", err)
		os.Exit(1)
	}
}
