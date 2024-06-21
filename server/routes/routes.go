package routes

import (
	"io/fs"
	"net/http"

	"github.com/abibby/remote-input/server/app/handlers"
	"github.com/abibby/remote-input/server/resources"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	// r.Use(auth.AttachUser())

	distContent, err := fs.Sub(resources.Content, "dist")
	if err != nil {
		panic(err)
	}
	r.Handle("/res", http.FileServerFS(distContent))

	r.Get("/", handlers.Home).Name("home")
	r.Get("/devices", handlers.BluetoothScanView).Name("devices")

	r.Group("/bluetooth", func(r *router.Router) {
		r.Get("/events", handlers.BluetoothEvents)
		r.Post("/scan/on", handlers.BluetoothScanOn)
		r.Post("/scan/off", handlers.BluetoothScanOff)
		r.Post("/device/{address}/connect", handlers.BluetoothConnect)
		r.Post("/device/{address}/pair", handlers.BluetoothPair)
		r.Post("/device/{address}/forget", handlers.BluetoothForget)
	})
}
