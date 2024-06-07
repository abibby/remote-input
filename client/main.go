package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/windows"
	"github.com/stephen-fox/user32util"
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

	lastMouseEvent := &common.MouseInputEvent{}
	b := make([]byte, 32)
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

		if e.Keyboard != nil {
			handleKeyboard(e.Keyboard)
		} else if e.Mouse != nil {
			handleMouse(e.Mouse, lastMouseEvent)
			lastMouseEvent = e.Mouse
		}

	}

}

func handleMouse(e *common.MouseInputEvent, last *common.MouseInputEvent) {
	var flags uint32
	if e.ButtonLeft() && !last.ButtonLeft() {
		flags |= user32util.MouseEventFLeftDown
	}
	if !e.ButtonLeft() && last.ButtonLeft() {
		flags |= user32util.MouseEventFLeftUp
	}
	if e.ButtonRight() && !last.ButtonRight() {
		flags |= user32util.MouseEventFRightDown
	}
	if !e.ButtonRight() && last.ButtonRight() {
		flags |= user32util.MouseEventFRightUp
	}

	if e.X != 0 || e.Y != 0 {
		flags |= user32util.MouseEventFMove
	}
	fmt.Printf("%08s\n", strconv.FormatInt(int64(e.Button), 2))
	windows.SendMouseInput(int32(e.X), int32(e.Y)*-1, flags)
}

func handleKeyboard(e *common.KeyboardInputEvent) {
	if e.EventType != common.EV_KEY {
		return
	}

	vKey, ok := keyMap[e.Code]
	if !ok {
		log.Printf("no map for key code %d", e.Code)
		return
	}
	var flag windows.KeyEventFlag
	if e.Value == 1 {
		flag = windows.KEYEVENTF_KEYPRESS
	} else if e.Value == 0 {
		flag = windows.KEYEVENTF_KEYUP
	}
	err := windows.SendInput(vKey, flag)
	if err != nil {
		log.Printf("failed to send input %v", err)
		return
	}

}
