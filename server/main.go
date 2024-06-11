package main

import (
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/abibby/remote-input/common"
)

var enabledEventTypes = [0x1f]bool{
	common.EV_KEY: true,
	common.EV_REL: true,
	common.EV_ABS: true,
}

func main() {

	dirById := "/dev/input/by-id"
	devicesById, err := os.ReadDir(dirById)
	if err != nil {
		log.Fatal(err)
	}

	devices := []string{}
	for _, f := range devicesById {
		devices = append(devices, path.Join(dirById, f.Name()))
	}

	listener, err := net.Listen("tcp", ":38808")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Print("listening")

	mux := NewConnMux()
	defer mux.Close()

	for _, device := range devices {
		if strings.HasSuffix(device, "-kbd") || strings.HasSuffix(device, "-event-mouse") || strings.HasSuffix(device, "-event-joystick") {
			go func(device string) {
				err = readDevice(device, mux)
				if err != nil {
					log.Printf("device %s failed: %v", device, err)
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

func readDevice(devicePath string, w io.Writer) error {
	f, err := os.Open(devicePath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("connected to %s", devicePath)

	e := common.InputEvent{}
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

		if !enabledEventTypes[e.EventType] {
			continue
		}
		// spew.Dump(e)

		out, err := e.MarshalBinary()
		if err != nil {
			log.Print(err)
			continue
		}
		_, err = w.Write(out)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}
