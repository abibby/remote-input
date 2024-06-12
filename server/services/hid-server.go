package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/server/config"
	"github.com/abibby/salusa/kernel"
)

// var enabledEventTypes = [0x1f]bool{
// 	common.EV_KEY: true,
// 	common.EV_REL: true,
// 	common.EV_ABS: true,
// }

type Device struct {
	Name  string
	Path  string
	Index uint16
	Type  int32
}

const eventPathBase = "/dev/input/event"
const dirById = "/dev/input/by-id"

type HIDServer struct {
	Cfg *config.Config `inject:""`
	Log *slog.Logger   `inject:""`
}

func NewHIDServer() *HIDServer {
	return &HIDServer{}
}

var _ kernel.Service = (*HIDServer)(nil)

func (h *HIDServer) Name() string {
	return "remote-input:hid-server"
}

// Run implements kernel.Service.
func (h *HIDServer) Run(ctx context.Context) error {
	devicesById, err := os.ReadDir(dirById)
	if err != nil {
		return (err)
	}

	devices := []*Device{}
	for _, f := range devicesById {
		p, err := filepath.EvalSymlinks(path.Join(dirById, f.Name()))
		if err != nil {
			h.Log.Warn("Failed to load device symlink", "error", err)
			continue
		}

		if !strings.HasPrefix(p, eventPathBase) {
			continue
		}

		strIndex := p[len(eventPathBase):]
		index, err := strconv.Atoi(strIndex)
		if err != nil {
			h.Log.Warn("Device path not in the expected form '/dev/input/event##'", "error", err)
			continue
		}
		deviceType := int32(-1)
		if strings.HasSuffix(f.Name(), "-kbd") {
			deviceType = common.DeviceTypeKeyboard
		} else if strings.HasSuffix(f.Name(), "-event-mouse") {
			deviceType = common.DeviceTypeMouse
		} else if strings.HasSuffix(f.Name(), "-event-joystick") {
			deviceType = common.DeviceTypeJoystick
		}
		devices = append(devices, &Device{
			Name:  f.Name(),
			Path:  p,
			Index: uint16(index),
			Type:  deviceType,
		})
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", h.Cfg.HIDPort))
	if err != nil {
		return (err)
	}
	defer listener.Close()

	h.Log.Info("Listening")

	mux := NewConnMux()
	defer mux.Close()

	for _, device := range devices {
		if device.Type != -1 {
			go func(device *Device) {
				err = h.readDevice(device, mux)
				if err != nil {
					h.Log.Warn("Device reading failed", "device", device.Name, "error", err)
				}
			}(device)
		}
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			h.Log.Warn("Error accepting new connection", "error", err)
			continue
		}
		mux.Add(conn)
	}
}

func (h *HIDServer) readDevice(device *Device, w io.Writer) error {
	f, err := os.Open(device.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	h.Log.Info("Device connected", "device", device.Name)

	e := common.InputEvent{}
	outBuffer := make([]byte, 0, e.Size()*8)
	b := make([]byte, e.Size())
	for {
		_, err = f.Read(b)
		if err != nil {
			return err
		}

		err = e.UnmarshalBinary(b)
		if err != nil {
			h.Log.Warn("Error decoding device event", "error", err)
			continue
		}

		// if !enabledEventTypes[e.EventType] {
		// 	continue
		// }
		if e.EventType == common.EV_SYN {
			e.Code = device.Index
			e.Value = device.Type
		}
		// spew.Dump(e)
		out, err := e.MarshalBinary()
		if err != nil {
			h.Log.Warn("Error encoding device event", "error", err)
			continue
		}
		outBuffer = append(outBuffer, out...)
		if e.EventType == common.EV_SYN {
			_, err = w.Write(outBuffer)
			if err != nil {
				h.Log.Warn("Error sending device events", "error", err)
				continue
			}
			outBuffer = outBuffer[:0]
		}
	}
}
