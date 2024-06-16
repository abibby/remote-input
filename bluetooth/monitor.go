package bluetooth

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
)

type Event interface {
}

type RemoveEvent struct {
	DevicePath string
}
type ChangeEvent struct{}
type ErrorEvent struct {
	Err error
}

func newEvent(s *dbus.Signal) (Event, error) {
	switch s.Name {
	case "org.freedesktop.DBus.ObjectManager.InterfacesRemoved":
		return newInterfaceRemovedEvent(s)
	case "org.freedesktop.DBus.Properties.PropertiesChanged":
		return newPropertiesChangedEvent(s)
	default:
		return nil, fmt.Errorf("")
	}
}

func newInterfaceRemovedEvent(s *dbus.Signal) (*RemoveEvent, error) {
	return &RemoveEvent{
		DevicePath: s.Body[0].(string),
	}, nil
}
func newPropertiesChangedEvent(s *dbus.Signal) (*ChangeEvent, error) {
	return &ChangeEvent{}, nil

}

func (bt *BT) Events(ctx context.Context) (chan Event, error) {

	remove, err := bt.addMatchSignal(dbus.WithMatchSender("org.bluez"))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to bluez: %w", err)
	}

	signals := make(chan *dbus.Signal, 10)
	events := make(chan Event, 10)
	bt.conn.Signal(signals)
	go func() {
		for s := range signals {
			if s.Sender != "org.bluez" {
				continue
			}

			e, err := newEvent(s)
			if err != nil {
				events <- &ErrorEvent{Err: err}
				continue
			}
			events <- e
		}
	}()

	go func() {
		<-ctx.Done()
		err = remove()
		if err != nil {
			panic(err)
		}
		close(signals)
	}()
	// a := &dbus.Signal{
	// 	Sender: ":1.7",
	// 	Path:   "/",
	// 	Name:   "org.freedesktop.DBus.ObjectManager.InterfacesRemoved",
	// 	Body: []interface{}{
	// 		"/org/bluez/hci0/dev_84_FC_E6_38_49_D6",
	// 		[]string{"org.freedesktop.DBus.Properties", "org.freedesktop.DBus.Introspectable", "org.bluez.Device1"},
	// 	},
	// 	Sequence: 0x14,
	// }
	// b := &dbus.Signal{
	// 	Sender: ":1.7",
	// 	Path:   "/org/bluez/hci0",
	// 	Name:   "org.freedesktop.DBus.Properties.PropertiesChanged",
	// 	Body: []interface{}{
	// 		"org.bluez.Adapter1",
	// 		map[string]dbus.Variant{"Discovering": dbus.Variant{sig: dbus.Signature{str: "b"}, value: false}},
	// 		[]string{},
	// 	},
	// 	Sequence: 0xf,
	// }
	return events, nil
}
