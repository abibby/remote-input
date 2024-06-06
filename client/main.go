package main

import (
	"log"
	"net"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/windows"
)

func main() {
	log.Printf("started")

	// serverIP := "192.168.2.50:38808"
	serverIP := "localhost:38808"

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("connected to %s", serverIP)

	b := make([]byte, 24)
	for {
		_, err = conn.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		e := &common.InputEvent{}

		err = e.UnmarshalBinary(b)
		if err != nil {
			log.Fatal(err)
		}

		if e.EventType == common.EV_KEY {
			vKey, ok := keyMap[e.Code]
			if !ok {
				log.Printf("no map for key code %d", e.Code)
				continue
			}
			if e.Value == 1 {
				windows.SendInput(vKey, windows.KEYEVENTF_KEYPRESS)
			} else if e.Value == 0 {
				windows.SendInput(vKey, windows.KEYEVENTF_KEYUP)
			}
		}
	}

}
