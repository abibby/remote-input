package main

import (
	"log"
	"net"

	"github.com/abibby/remote-input/common"
)

func main() {
	listener, err := net.Listen("tcp", ":38808")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Print("listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go serve(conn)

	}
}

func serve(conn net.Conn) {
	defer conn.Close()

	var err error
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

		log.Printf("%v\n", e)
	}
}
