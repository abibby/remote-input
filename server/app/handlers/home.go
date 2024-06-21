package handlers

import (
	"github.com/abibby/remote-input/server/services"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/view"
)

type HomeRequest struct{}

var Home = request.Handler(func(r *HomeRequest) (*view.ViewHandler, error) {
	return view.View("index.html", map[string]any{
		"Devices": services.ConnectedDevices,
	}), nil
})
