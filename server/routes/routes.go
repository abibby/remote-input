package routes

import (
	"io/fs"
	"net/http"

	"github.com/abibby/remote-input/server/app/handlers"
	"github.com/abibby/remote-input/server/resources"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/view"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	// r.Use(auth.AttachUser())

	distContent, err := fs.Sub(resources.Content, "dist")
	if err != nil {
		panic(err)
	}
	r.Handle("/res", http.FileServerFS(distContent))

	r.Get("/", view.View("index.html", nil)).Name("home")
	r.Get("/scan", view.View("scan.html", nil)).Name("scan")

	r.Group("/bluetooth", func(r *router.Router) {
		r.Get("/scan", handlers.BluetoothScan)
		r.Post("/connect/{address}", handlers.BluetoothConnect)
	})
}
