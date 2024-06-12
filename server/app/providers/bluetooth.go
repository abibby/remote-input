package providers

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/di"
	"tinygo.org/x/bluetooth"
)

func RegisterBluetoothAdapter(ctx context.Context) error {
	di.RegisterLazySingleton(ctx, func() (*bluetooth.Adapter, error) {
		adapter := bluetooth.DefaultAdapter
		err := adapter.Enable()
		if err != nil {
			return nil, fmt.Errorf("could not enable bluetooth adapter: %w", err)
		}
		return adapter, nil
	})
	return nil
}
