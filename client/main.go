package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/config"
	"github.com/abibby/remote-input/windows"
)

func main() {
	log.Printf("started")

	serverIP := fmt.Sprintf("%s:%d", config.Host, config.Port)

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("connected to %s", serverIP)

	events := make([]common.InputEvent, 0, 8)

	b := make([]byte, 24)
	for {
		events = append(events, common.InputEvent{})
		e := &events[len(events)-1]

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

		if e.EventType == common.EV_SYN {
			err = handleEvent(events)
			if err != nil {
				log.Print(err)
			}
			events = events[:0]
		}
	}

}

func handleEvent(events []common.InputEvent) error {
	syn := events[len(events)-1]
	rest := events[:len(events)-1]
	// if syn.Value == common.DeviceTypeKeyboard {
	// 	return handleKeyboard(rest)
	// }
	// if syn.Value == common.DeviceTypeMouse {
	// 	return handleMouse(rest)
	// }
	if syn.Value == common.DeviceTypeJoystick {
		return handleJoystick(rest, syn.Code)
	}

	return nil
}

func handleJoystick(events []common.InputEvent, index uint16) error {
	keyboardEvents := []common.InputEvent{}

	for _, e := range events {
		log.Printf("%v: %x %x\n", e.EventType, e.Code, e.Value)
		switch e.EventType {
		case common.EV_ABS:
		case common.EV_KEY:
			if e.Code < uint16(len(keyMap)) {
				keyboardEvents = append(keyboardEvents, e)
			} else {
				buttonID := e.Code - common.JOYSTICK_BASE
				log.Printf("gamepad button %d\n", buttonID)
			}
		}
	}
	fmt.Println()
	return nil
}
func handleMouse(events []common.InputEvent) error {
	var flags uint32
	var data int32
	var dx int32
	var dy int32
	for _, e := range events {
		switch e.EventType {
		case common.EV_REL:
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
			}
		case common.EV_KEY:
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
		}
	}
	return windows.SendMouseInput(dx, dy, data, flags)
}
func handleKeyboard(events []common.InputEvent) error {
	for _, e := range events {
		if e.EventType != common.EV_KEY {
			continue
		}
		vKey := keyMap[e.Code]
		if vKey == 0 {
			return fmt.Errorf("no map for key code %d", e.Code)
		}

		var flag windows.KeyEventFlag
		if e.Value == 0 {
			flag = windows.KEYEVENTF_KEYUP
		} else if e.Value == 1 {
			flag = windows.KEYEVENTF_KEYPRESS
		} else if e.Value == 2 {
			flag = windows.KEYEVENTF_KEYPRESS
		}
		return windows.SendInput(vKey, flag)
	}
	return nil
}
