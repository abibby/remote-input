package providers

import (
	"context"
	"fmt"
	"sync"

	"github.com/abibby/remote-input/server/config"
	"github.com/abibby/salusa/di"
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
)

type Agent struct {
	agent.SimpleAgent

	mtx      *sync.RWMutex
	pincodes []PinCodeEvents
}

func (a *Agent) DisplayPinCode(device dbus.ObjectPath, pincode string) *dbus.Error {
	a.mtx.RLock()
	defer a.mtx.Unlock()

	for _, p := range a.pincodes {
		p <- pincode
	}
	return nil
}
func (a *Agent) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error {
	return nil
}

type PinCodeEvents chan string

func RegisterBluetoothAdapter(ctx context.Context) error {
	// log.SetLevel(log.ErrorLevel)
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	ag := &Agent{
		SimpleAgent: *agent.NewDefaultSimpleAgent(),
		mtx:         &sync.RWMutex{},
		pincodes:    []PinCodeEvents{},
	}
	err = agent.ExposeAgent(conn, ag, agent.CapKeyboardOnly, true)
	if err != nil {
		return fmt.Errorf("SimpleAgent: %s", err)
	}

	di.RegisterSingleton(ctx, func() *dbus.Conn {
		return conn
	})
	di.RegisterSingleton(ctx, func() *Agent {
		return ag
	})
	di.RegisterSingleton(ctx, func() agent.Agent1Client {
		return ag
	})
	di.RegisterLazySingletonWith(ctx, func(cfg *config.Config) (*adapter.Adapter1, error) {
		a, err := adapter.GetAdapter(cfg.AdapterID)
		if err != nil {
			return nil, err
		}
		return a, nil
	})

	di.RegisterWith(ctx, func(ctx context.Context, tag string, agent *Agent) (PinCodeEvents, error) {
		agent.mtx.Lock()
		defer agent.mtx.Unlock()

		pincodes := make(PinCodeEvents)
		agent.pincodes = append(agent.pincodes, pincodes)
		return pincodes, nil
	})

	return nil
}
