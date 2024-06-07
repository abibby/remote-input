package main

import (
	"io"
	"log"
	"net"

	"github.com/abibby/remote-input/common"
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

	e := &common.InputEvent{}
	b := make([]byte, 32)
	for {
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

		// switch e := e.(type) {
		// case *common.KeyboardInputEvent:
		// 	if e.EventType == common.EV_KEY {
		// 		vKey, ok := keyMap[e.Code]
		// 		if !ok {
		// 			log.Printf("no map for key code %d", e.Code)
		// 			continue
		// 		}
		// 		var flag windows.KeyEventFlag
		// 		if e.Value == 1 {
		// 			flag = windows.KEYEVENTF_KEYPRESS
		// 		} else if e.Value == 0 {
		// 			flag = windows.KEYEVENTF_KEYUP
		// 		}
		// 		err = windows.SendInput(vKey, flag)
		// 		if err != nil {
		// 			log.Printf("failed to send input %v", err)
		// 			continue
		// 		}
		// 	}
		// case *common.MouseInputEvent:
		// }
		log.Printf("%#v\n", e)
		// e := &common.KeyboardInputEvent{}

		// err = e.UnmarshalBinary(b)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// if e.EventType == common.EV_KEY {
		// 	vKey, ok := keyMap[e.Code]
		// 	if !ok {
		// 		log.Printf("no map for key code %d", e.Code)
		// 		continue
		// 	}
		// 	if e.Value == 1 {
		// 		windows.SendInput(vKey, windows.KEYEVENTF_KEYPRESS)
		// 	} else if e.Value == 0 {
		// 		windows.SendInput(vKey, windows.KEYEVENTF_KEYUP)
		// 	}
		// }
	}

}
