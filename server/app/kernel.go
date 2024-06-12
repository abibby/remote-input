package app

import (
	"github.com/abibby/remote-input/server/app/models"
	"github.com/abibby/remote-input/server/app/providers"
	"github.com/abibby/remote-input/server/config"
	"github.com/abibby/remote-input/server/migrations"
	"github.com/abibby/remote-input/server/resources"
	"github.com/abibby/remote-input/server/routes"
	"github.com/abibby/remote-input/server/services"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/salusadi"
	"github.com/abibby/salusa/view"
)

var Kernel = kernel.New[*config.Config](
	kernel.Config(config.Load),
	kernel.Bootstrap(
		salusadi.Register[*models.User](migrations.Use()),
		view.Register(resources.Content, "**/*.html"),
		providers.Register,
		providers.RegisterBluetoothAdapter,
	),
	kernel.Services(
		// cron.Service().
		// 	Schedule("* * * * *", &events.LogEvent{Message: "cron event"}),
		// event.Service(
		// 	event.NewListener[*jobs.LogJob](),
		// ),
		services.NewHIDServer(),
	),
	router.InitRoutes(routes.InitRoutes),
)
