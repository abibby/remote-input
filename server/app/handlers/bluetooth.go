package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/abibby/remote-input/server/eventsource"
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

		scanResults := make(chan bluetooth.ScanResult)

		go func() {
			err := r.Adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) {
				scanResults <- sr
			})
			if err != nil {
				r.Log.Error("Scan failed", "error", err)
				return
			}
		}()

		defer func() {
			err := r.Adapter.StopScan()
			if err != nil {
				r.Log.Error("Stop scan failed", "error", err)
				return
			}
			r.Log.Info("scan stopped")
		}()
		added := set.New[string]()
		es := eventsource.New(w)
		for {
			select {
			case scanResult := <-scanResults:
				addr := scanResult.Address.String()
				if added.Has(addr) || scanResult.LocalName() == "" {
					continue
				}
				added.Add(addr)

				b, err := view.View("scan-result.html", &Device{
					Address:   addr,
					RSSI:      scanResult.RSSI,
					LocalName: scanResult.LocalName(),
				}).Bytes(r.Ctx)
				if err != nil {
					r.Log.Warn("failed to load device view", "error", err)
					return
				}
				err = es.Send(&eventsource.Event{
					Event: "scan-results",
					Data:  string(b),
				})
				if err != nil {
					r.Log.Warn("failed to send scan results", "error", err)
					return
				}
			case <-r.Ctx.Done():
				r.Log.Info("closed")
				return
			}
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
