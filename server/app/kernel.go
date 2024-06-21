package app

import (
	"github.com/abibby/remote-input/server/app/providers"
	"github.com/abibby/remote-input/server/config"
	"github.com/abibby/remote-input/server/resources"
	"github.com/abibby/remote-input/server/routes"
	"github.com/abibby/remote-input/server/services"
	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/filesystem"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/view"
)

var Kernel = kernel.New[*config.Config](
	kernel.Config(config.Load),
	kernel.Bootstrap(
		clog.Register(nil),
		request.Register,
		filesystem.Register,
		view.Register(resources.Content, "**/*.html"),
		providers.Register,
		providers.RegisterBluetoothAdapter,
	),
	kernel.Services(
		services.NewHIDServer(),
	),
	router.InitRoutes(routes.InitRoutes),
)
