package services

import (
	"bytes"
	"context"
	"encoding"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/server/config"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/set"
)

var ConnectedDevices = []string{}

type Device struct {
	Name     string
	Phys     string
	Sysfs    string
	Uniq     string
	Handlers []string

	Path  string
	Index uint16
	Type  int32
}

var _ encoding.TextUnmarshaler = (*Device)(nil)

var eventRE = regexp.MustCompile(`event\d+`)
var joystickRE = regexp.MustCompile(`js\d+`)
var mouseRE = regexp.MustCompile(`mouse\d+`)

func (d *Device) UnmarshalText(text []byte) error {
	lines := bytes.Split(text, []byte("\n"))
	d.Type = -1
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		_, key, value := splitLine(line)
		switch key {
		case "Name":
			d.Name = strings.Trim(value, `"`)
		case "Phys":
			d.Phys = value
		case "Sysfs":
			d.Sysfs = value
		case "Uniq":
			d.Uniq = value
		case "Handlers":
			d.Handlers = strings.Split(value, " ")

			event := eventRE.FindString(value)

			d.Path = "/dev/input/" + event
			if joystickRE.MatchString(value) {
				d.Type = common.DeviceTypeJoystick
			} else if mouseRE.MatchString(value) {
				d.Type = common.DeviceTypeMouse
			} else if strings.Contains(value, "kbd") {
				d.Type = common.DeviceTypeKeyboard
			}
		}
	}
	return nil
}

func splitLine(line []byte) (byte, string, string) {
	parts := bytes.SplitN(line, []byte{'='}, 2)
	start := parts[0]
	value := parts[1]

	return start[0], string(start[3:]), string(value)
}

// func (d *Device) UnmarshalBinary()

// const eventPathBase = "/dev/input/event"
// const dirById = "/dev/input/by-id"

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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", h.Cfg.HIDPort))
	if err != nil {
		return (err)
	}
	defer listener.Close()

	h.Log.Info("Listening")

	mux := NewConnMux()
	defer mux.Close()

	go func() {
		ticker := time.NewTicker(time.Second)
		devices := map[string]*Device{}
		cancels := map[string]context.CancelFunc{}
		for range ticker.C {
			hasChange := false
			found := set.New[string]()
			newDevices, err := h.getDevices()
			if err != nil {
				h.Log.Info("failed to get devices", "error", err)
				continue
			}
			for _, device := range newDevices {
				path := device.Path
				_, ok := devices[path]
				if !ok {
					h.Log.Info("device", "name", device.Name, "ok", ok)
					ctx, cancel := context.WithCancel(ctx)
					devices[path] = device
					cancels[path] = cancel

					hasChange = true
					go h.readDevice(ctx, device, mux)
				}
				found.Add(path)
			}
			for path := range devices {
				if found.Has(path) {
					continue
				}
				hasChange = true
				cancels[path]()
				delete(devices, path)
				delete(cancels, path)
			}
			if hasChange {
				ConnectedDevices = make([]string, len(newDevices))
				for i, device := range newDevices {
					ConnectedDevices[i] = device.Name
				}
			}
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			h.Log.Warn("Error accepting new connection", "error", err)
			continue
		}
		mux.Add(conn)
	}
}

func (h *HIDServer) readDevice(ctx context.Context, device *Device, w io.Writer) {
	f, err := os.Open(device.Path)
	if err != nil {
		h.Log.Warn("Device open failed", "device", device.Name, "error", err)
		return
	}
	defer f.Close()

	h.Log.Info("Device connected", "device", device.Name)

	e := common.InputEvent{}
	outBuffer := make([]byte, 0, e.Size()*8)
	b := make([]byte, e.Size())
	for ctx.Err() == nil {
		_, err = f.Read(b)
		if err != nil {
			if errors.Is(err, syscall.Errno(0x13)) {
				h.Log.Warn("Device removed", "device", device.Name)
				return
			}
			h.Log.Warn("Device reading failed", "device", device.Name, "error", err)
			continue
		}

		err = e.UnmarshalBinary(b)
		if err != nil {
			h.Log.Warn("Error decoding device event", "error", err)
			continue
		}

		if e.EventType == common.EV_SYN {
			e.Code = device.Index
			e.Value = device.Type
		}

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

	if ctx.Err() != nil {
		h.Log.Error("context closed", "error", ctx.Err())
	}
}

func (h *HIDServer) getDevices() ([]*Device, error) {

	b, err := os.ReadFile("/proc/bus/input/devices")
	if err != nil {
		return nil, err
	}

	bDevices := bytes.Split(b, []byte("\n\n"))

	devices := make([]*Device, 0, len(bDevices))
	for _, b := range bDevices {
		d := &Device{}
		err = d.UnmarshalText(b)
		if err != nil {
			return nil, err
		}
		if d.Type == -1 {
			continue
		}
		devices = append(devices, d)
	}
	return devices, nil
}
