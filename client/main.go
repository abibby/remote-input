package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/windows"
)

func main() {
	log.Printf("started")

	// serverIP := os.Getenv("REMOTE_INPUT_HOST")
	// serverIP := "192.168.2.54:38808"
	// serverIP := "192.168.2.38:38808"
	serverIP := "localhost:38808"

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
		case common.EV_ABS:
			err = handleAbs(e)
		default:
			log.Printf("unhandled event type %v\n", e.EventType)
		}
		if err != nil {
			log.Print(err)
			continue
		}
	}

}

func handleAbs(e *common.InputEvent) error {
	return nil
}
func handleRel(e *common.InputEvent) error {
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
		flags |= MouseEventFWheel
	default:
		// fmt.Printf("% 4d % 4d\n", e.Code, e.Value)
		return nil
	}
	return windows.SendMouseInput(dx, dy, data, flags)
}

func handleKey(e *common.InputEvent) error {
	if e.Code >= common.JOYSTICK_BASE {
		// Joystick
		buttonNum := e.Code - common.JOYSTICK_BASE
		log.Printf("joystick button %d", buttonNum)
		return nil
	} else if e.Code > uint16(len(keyMap)) {
		// Mouse
		var flags uint32
		switch e.Code {
		case common.MOUSE_LEFT:
			if e.Value == 1 {
				flags |= MouseEventFLeftDown
			} else if e.Value == 0 {
				flags |= MouseEventFLeftUp
			}
		case common.MOUSE_RIGHT:
			if e.Value == 1 {
				flags |= MouseEventFRightDown
			} else if e.Value == 0 {
				flags |= MouseEventFRightUp
			}
		}
		return windows.SendMouseInput(0, 0, 0, flags)
	}

	vKey := keyMap[e.Code]
	if vKey == 0 {
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
