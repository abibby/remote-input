package main

import (
	"log"
	"net"
	"os"
)

func main() {
	log.Printf("started")

	dev := "/dev/input/by-id/usb-Generic_USB_Keyboard-event-kbd"
	// serverIP := "192.168.2.50:38808"
	serverIP := "localhost:38808"

	f, err := os.Open(dev)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Printf("connected to %s", dev)

	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to %s", serverIP)

	b := make([]byte, 24)
	for {
		_, err = f.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%v", len(b))

		// e := &common.InputEvent{}

		// err = e.UnmarshalBinary(b)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		_, err = conn.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}

}
