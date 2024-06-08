package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/windows"
)

func main() {
	log.Printf("started")

	serverIP := "192.168.2.38:38808"
	// serverIP := "localhost:38808"

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("connected to %s", serverIP)

	b := make([]byte, 24)
	for {
		e := &common.InputEvent{}
		_, err := conn.Read(b)
		if err == io.EOF {
			log.Print("disconnected")
			return
		} else if err != nil {
			log.Print(err)
			continue
		}

		err = e.UnmarshalBinary(b)
		if err != nil {
			log.Print(err)
			continue
		}

		switch e.EventType {
		case common.EV_KEY:
			err = handleKey(e)
		case common.EV_REL:
			err = handleRel(e)
		}
		if err != nil {
			log.Print(err)
			continue
		}
	}

}

var lastMouseTime time.Time

func handleRel(e *common.InputEvent) error {
	if !lastMouseTime.Equal(e.Time) {
		lastMouseTime = e.Time

		fmt.Println()
	}
	var flags uint32
	var data int32
	var dx int32
	var dy int32
	switch e.Code {
	case 0:
		dx = e.Value
		flags |= MouseEventFMove
	case 1:
		dy = e.Value
		flags |= MouseEventFMove
	case 11:
		data = e.Value
		flags |= MouseEventFHWheel
	default:
		fmt.Printf("% 4d % 4d\n", e.Code, e.Value)
	}
	return windows.SendMouseInput(dx, dy*-1, data, flags)
}

func handleKey(e *common.InputEvent) error {
	vKey, ok := keyMap[e.Code]
	if !ok {
		return fmt.Errorf("no map for key code %d", e.Code)
	}
	var flag windows.KeyEventFlag
	if e.Value == 1 {
		flag = windows.KEYEVENTF_KEYPRESS
	} else if e.Value == 0 {
		flag = windows.KEYEVENTF_KEYUP
	}
	return windows.SendInput(vKey, flag)
}
