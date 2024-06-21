package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abibby/remote-input/server/app"
	"github.com/abibby/salusa/di"
	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	ctx := di.ContextWithDependencyProvider(
		context.Background(),
		di.NewDependencyProvider(),
	)

	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		return fmt.Errorf("error bootstrapping: %w", err)
	}

	go p.run(ctx)
	return nil
}
func (p *program) run(ctx context.Context) {
	err := app.Kernel.Run(ctx)
	if err != nil {
		app.Kernel.Logger(ctx).Error("error running", "error", err)
		os.Exit(1)
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "remote-input-server",
		DisplayName: "Remote Input Server",
		Description: "Server to send input events to a client.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err = s.Install()
		case "uninstall":
			err = s.Uninstall()
		default:
			err = fmt.Errorf("unknown command %s, expected install or uninstall", os.Args[1])
		}
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
