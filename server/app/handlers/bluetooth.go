package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/set"
	"github.com/abibby/salusa/view"
	"tinygo.org/x/bluetooth"
)

type BluetoothScanRequest struct {
	Adapter *bluetooth.Adapter `inject:""`
	Log     *slog.Logger       `inject:""`
	Ctx     context.Context    `inject:""`
}

type Device struct {
	Address   string
	RSSI      int16
	LocalName string
}

var BluetoothScan = request.Handler(func(r *BluetoothScanRequest) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		added := set.New[string]()
		go func() {
			err := r.Adapter.Scan(func(a *bluetooth.Adapter, device bluetooth.ScanResult) {
				addr := device.Address.String()
				if added.Has(addr) || device.LocalName() == "" {
					return
				}
				added.Add(addr)

				err := view.View("scan-result.html", &Device{
					Address:   addr,
					RSSI:      device.RSSI,
					LocalName: device.LocalName(),
				}).Execute(r.Ctx, w)
				if err != nil {
					r.Log.Warn("failed to load device view", "error", err)
					return
				}
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			})
			if err != nil {
				r.Log.Error("Scan failed", "error", err)
				return
			}
		}()

		select {
		case <-r.Ctx.Done():
			r.Log.Info("closed")
		case <-time.After(time.Minute * 5):
			r.Log.Info("timeout")
		}

		err := r.Adapter.StopScan()
		if err != nil {
			r.Log.Error("Stop scan failed", "error", err)
			return
		}
	}), nil
})

type BluetoothConnectRequest struct {
	Address string             `path:"address"`
	Adapter *bluetooth.Adapter `inject:""`
	Log     *slog.Logger       `inject:""`
	Ctx     context.Context    `inject:""`
}

var BluetoothConnect = request.Handler(func(r *BluetoothConnectRequest) (http.Handler, error) {
	address, err := ParseAddress(r.Address)
	if err != nil {
		return nil, err
	}
	dev, err := r.Adapter.Connect(address, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, err
	}
	return view.View("connected.html", &Device{
		Address: dev.Address.String(),
	}), nil
})
