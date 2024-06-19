package providers

import (
	"context"
	"fmt"

	"github.com/abibby/remote-input/server/config"
	"github.com/abibby/salusa/di"
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	log "github.com/sirupsen/logrus"
)

func RegisterBluetoothAdapter(ctx context.Context) error {
	log.SetLevel(log.ErrorLevel)
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	ag := agent.NewSimpleAgent()
	err = agent.ExposeAgent(conn, ag, agent.CapKeyboardOnly, true)
	if err != nil {
		return fmt.Errorf("SimpleAgent: %s", err)
	}

	di.RegisterSingleton(ctx, func() *dbus.Conn {
		return conn
	})
	di.RegisterLazySingletonWith(ctx, func(cfg *config.Config) (*adapter.Adapter1, error) {
		a, err := adapter.GetAdapter(cfg.AdapterID)
		if err != nil {
			return nil, err
		}
		return a, nil
	})

	return nil
}
