package routes

import (
	"github.com/abibby/remote-input/server/app/handlers"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	// r.Use(auth.AttachUser())

	// auth.RegisterRoutes(r, func(r *auth.EmailVerifiedUserCreateRequest) *models.User {
	// 	return &models.User{
	// 		EmailVerifiedUser: *auth.NewEmailVerifiedUser(r),
	// 	}
	// }, "/reset-password")

	// r.Get("/", view.View("index.html", nil)).Name("home")
	// r.Get("/login", view.View("login.html", nil)).Name("login")
	// r.Get("/user/create", view.View("create_user.html", nil)).Name("user.create")

	// r.Get("/user", handlers.UserList)
	// r.Get("/user/{id}", handlers.UserGet)

	r.Group("/bluetooth", func(r *router.Router) {
		r.Get("/scan", handlers.BluetoothScan)
	})
}
