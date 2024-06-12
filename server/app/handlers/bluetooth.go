package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/abibby/salusa/request"
	"tinygo.org/x/bluetooth"
)

type BluetoothScanRequest struct {
	Adapter *bluetooth.Adapter `inject:""`
	Log     *slog.Logger       `inject:""`
}

var BluetoothScan = request.Handler(func(r *BluetoothScanRequest) (http.Handler, error) {

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		jw := json.NewEncoder(w)
		r.Log.Info("starting scan")
		err := r.Adapter.Scan(func(a *bluetooth.Adapter, device bluetooth.ScanResult) {
			jw.Encode(device)
			println("found device:", device.Address.String(), device.RSSI, device.LocalName())

		})
		if err != nil {
			r.Log.Error("Scan failed", "error", err)
			return
		}
		r.Log.Info("its async")
		defer func() {
			err = r.Adapter.StopScan()
			if err != nil {
				r.Log.Error("Stop scan failed", "error", err)
				return
			}
		}()

		time.Sleep(time.Second * 5)

	}), nil
})
