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

func main() {

	dir := "/dev/input/by-id"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	devices := make([]string, 0, len(files))
	for _, f := range files {
		devices = append(devices, path.Join(dir, f.Name()))
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
		if strings.HasSuffix(device, "-kbd") || strings.HasSuffix(device, "-event-mouse") {
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

		if e.EventType != common.EV_KEY && e.EventType != common.EV_REL {
			continue
		}

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
