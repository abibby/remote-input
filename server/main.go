package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abibby/remote-input/bluetoothctl"
	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/config"
	"github.com/davecgh/go-spew/spew"
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

func main() {

	devs, errs := bluetoothctl.Scan()
	for {
		select {
		case device := <-devs:
			spew.Dump(device)
		case err := <-errs:
			log.Fatal(err)
		}
	}

	spew.Dump("test")

	devicesById, err := os.ReadDir(dirById)
	if err != nil {
		log.Fatal(err)
	}

	devices := []*Device{}
	for _, f := range devicesById {
		p, err := filepath.EvalSymlinks(path.Join(dirById, f.Name()))
		if err != nil {
			log.Print(err)
			continue
		}

		if !strings.HasPrefix(p, eventPathBase) {
			continue
		}

		strIndex := p[len(eventPathBase):]
		index, err := strconv.Atoi(strIndex)
		if err != nil {
			log.Print(err)
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Print("listening")

	mux := NewConnMux()
	defer mux.Close()

	for _, device := range devices {
		if device.Type != -1 {
			go func(device *Device) {
				err = readDevice(device, mux)
				if err != nil {
					log.Printf("device %s failed: %v", device.Name, err)
				}
			}(device)
		}
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		mux.Add(conn)
	}
}

func readDevice(device *Device, w io.Writer) error {
	f, err := os.Open(device.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("connected to %s", device.Name)

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
			log.Print(err)
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
			log.Print(err)
			continue
		}
		outBuffer = append(outBuffer, out...)
		if e.EventType == common.EV_SYN {
			_, err = w.Write(outBuffer)
			if err != nil {
				log.Print(err)
				continue
			}
			outBuffer = outBuffer[:0]
		}
	}
}
