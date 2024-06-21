package handlers

import (
	"context"
	"log"
	"log/slog"
	"slices"

	"github.com/abibby/remote-input/server/app/providers"
	"github.com/abibby/remote-input/server/eventsource"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/view"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

type BluetoothScanViewRequest struct {
	Adapter *adapter.Adapter1 `inject:""`
}

var BluetoothScanView = request.Handler(func(r *BluetoothScanViewRequest) (*view.ViewHandler, error) {
	devices, err := r.Adapter.GetDevices()
	if err != nil {
		return nil, err
	}
	sortDevices(devices, nil)

	return view.View("scan.html", map[string]any{"Devices": devices}), nil
})

type BluetoothEventsRequest struct {
	Adapter       *adapter.Adapter1       `inject:""`
	PinCodeEvents providers.PinCodeEvents `inject:""`
	Log           *slog.Logger            `inject:""`
	Ctx           context.Context         `inject:""`
}

var BluetoothEvents = request.Handler(func(r *BluetoothEventsRequest) (*eventsource.EventSource, error) {
	omSignal, omCancel, err := r.Adapter.GetObjectManagerSignal()
	if err != nil {
		return nil, err
	}
	es := make(chan *eventsource.Event, 10)

	go func() {
		defer func() {
			omCancel()
			close(es)
		}()
		var lastDevices []*device.Device1
		for {
			select {
			case s := <-omSignal:
				if s == nil {
					return
				}
				log.Print(s.Name)

				if s.Name != bluez.InterfacesAdded && s.Name != bluez.InterfacesRemoved {
					r.Log.Info("unknown event", "name", s.Name)
					continue
				}

				devices, err := r.Adapter.GetDevices()
				if err != nil {
					r.Log.Warn("failed to get connected devices")
					continue
				}

				sortDevices(devices, lastDevices)
				lastDevices = devices

				b, err := view.View("scan-results.html", devices).Bytes(r.Ctx)
				if err != nil {
					r.Log.Warn("failed to build scan results html")
					continue
				}
				es <- &eventsource.Event{
					Event: "scan-results",
					Data:  string(b),
				}
			case pin := <-r.PinCodeEvents:
				es <- &eventsource.Event{
					Event: "display-pin-code",
					Data:  pin,
				}
			case <-r.Ctx.Done():
				return
			}
		}
	}()
	return eventsource.New(es).SetErrorHandler(func(err error) {
		r.Log.Warn("event source send error", "error", err)
	}), nil
})

func sortDevices(devices, lastDevices []*device.Device1) {
	lastIndex := map[string]int{}
	for i, d := range lastDevices {
		lastIndex[d.Properties.Address] = i
	}
	slices.SortFunc(devices, func(a, b *device.Device1) int {
		ai, aok := lastIndex[a.Properties.Address]
		bi, bok := lastIndex[b.Properties.Address]
		if aok && bok {
			return ai - bi
		}
		if aok {
			return -1
		}
		if bok {
			return 1
		}
		if a.Properties.Alias == b.Properties.Alias {
			return 0
		}
		if a.Properties.Alias > b.Properties.Alias {
			return 1
		}
		return -1
	})
}

type BluetoothScanOnRequest struct {
	Adapter *adapter.Adapter1 `inject:""`
}

var BluetoothScanOn = request.Handler(func(r *BluetoothScanOnRequest) (string, error) {
	err := r.Adapter.StartDiscovery()
	if err != nil {
		return "", err
	}
	err = r.Adapter.SetPairable(true)
	if err != nil {
		return "", err
	}
	return "Scanning...", nil
})

type BluetoothScanOffRequest struct {
	Adapter *adapter.Adapter1 `inject:""`
}

var BluetoothScanOff = request.Handler(func(r *BluetoothScanOffRequest) (string, error) {
	err := r.Adapter.SetPairable(false)
	if err != nil {
		return "", err
	}
	err = r.Adapter.StopDiscovery()
	if err != nil {
		return "", err
	}
	return "Scan stopped", nil
})

type BluetoothConnectRequest struct {
	Address string            `path:"address"`
	Adapter *adapter.Adapter1 `inject:""`
}

var BluetoothConnect = request.Handler(func(r *BluetoothConnectRequest) (string, error) {
	device, err := r.Adapter.GetDeviceByAddress(r.Address)
	if err != nil {
		return "", err
	}

	err = device.Connect()
	if err != nil {
		return "", err
	}

	return "connected", nil
})

type BluetoothPairRequest struct {
	Address string            `path:"address"`
	Adapter *adapter.Adapter1 `inject:""`
}

var BluetoothPair = request.Handler(func(r *BluetoothPairRequest) (string, error) {
	device, err := r.Adapter.GetDeviceByAddress(r.Address)
	if err != nil {
		return "", err
	}

	err = device.Pair()
	if err != nil {
		return "", err
	}

	err = device.SetTrusted(true)
	if err != nil {
		return "", err
	}

	err = device.Connect()
	if err != nil {
		return "", err
	}

	return "paired", nil
})

type BluetoothForgetRequest struct {
	Address string            `path:"address"`
	Adapter *adapter.Adapter1 `inject:""`
}

var BluetoothForget = request.Handler(func(r *BluetoothForgetRequest) (string, error) {
	device, err := r.Adapter.GetDeviceByAddress(r.Address)
	if err != nil {
		return "", err
	}

	err = r.Adapter.RemoveDevice(device.Path())
	if err != nil {
		return "", err
	}

	return "forget", nil
})
