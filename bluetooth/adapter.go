package bluetooth

import "github.com/godbus/dbus/v5"

// https://github.com/bluez/bluez/blob/master/doc/org.bluez.Adapter.rst
type Adapter struct {
	obj dbus.BusObject
}

func GetAdapter(conn *dbus.Conn, path string) *Adapter {
	obj := conn.Object("org.bluez.Adapter1", conn.BusObject().Path())
	return &Adapter{
		obj: obj,
	}
}

// StartDiscovery starts device discovery session which may include starting an
// inquiry and/or scanning procedures and remote device name resolving.
//
// Use StopDiscovery to release the sessions acquired.
//
// This process will start creating Device objects as new devices are
// discovered.
//
// During discovery RSSI delta-threshold is imposed.
//
// Each client can request a single device discovery session per adapter.
//
// Possible errors:
//
//   - org.bluez.Error.NotReady:
//   - org.bluez.Error.Failed:
//   - org.bluez.Error.InProgress:
func (a *Adapter) StartDiscovery() {
	c := a.obj.Call("StartDiscovery", 0)
	<-c.Done
}

// Stops device discovery session started by StartDiscovery.
//
// Note that a discovery procedure is shared between all discovery sessions thus
// calling StopDiscovery will only release a single session and discovery will
// stop when all sessions from all clients have finished.
//
// Possible errors:
//
//   - org.bluez.Error.NotReady:
//   - org.bluez.Error.Failed:
//   - org.bluez.Error.NotAuthorized:
func (a *Adapter) StopDiscovery() {
	c := a.obj.Call("StopDiscovery", 0)
	<-c.Done
}
