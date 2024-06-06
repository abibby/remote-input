package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/abibby/remote-input/common"
)

type ConnMux struct {
	mtx   *sync.Mutex
	conns []net.Conn
}

func NewConnMux() *ConnMux {
	return &ConnMux{
		mtx:   &sync.Mutex{},
		conns: []net.Conn{},
	}
}

// Write implements io.Writer.
func (m *ConnMux) Write(p []byte) (n int, err error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	maxN := -1
	errs := make([]error, 0, len(m.conns))
	toRemove := make([]net.Conn, 0, len(m.conns))
	for _, c := range m.conns {
		n, err = c.Write(p)
		if err != nil {
			toRemove = append(toRemove, c)
			if err != io.EOF {
				errs = append(errs, err)
			}
		}
		if n > maxN {
			maxN = n
		}
	}

	for _, conn := range toRemove {
		m.remove(conn)
	}

	if len(errs) > 0 {
		return 0, errors.Join(errs...)
	}

	return maxN, nil
}

// Close implements io.WriteCloser.
func (m *ConnMux) Close() error {
	errs := make([]error, 0, len(m.conns))
	for _, c := range m.conns {
		err := c.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (m *ConnMux) Add(conn net.Conn) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.add(conn)
}
func (m *ConnMux) add(conn net.Conn) {
	m.conns = append(m.conns, conn)
}
func (m *ConnMux) Remove(conn net.Conn) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.remove(conn)
}
func (m *ConnMux) remove(conn net.Conn) error {
	for i, c := range m.conns {
		if c == conn {
			if i != len(m.conns)-1 {
				m.conns[i] = m.conns[len(m.conns)-1]
			}
			m.conns = m.conns[:len(m.conns)-1]

			return c.Close()
		}
	}
	return nil
}

var _ io.WriteCloser = (*ConnMux)(nil)

func main() {
	dev := "/dev/input/by-id/usb-Generic_USB_Keyboard-event-kbd"

	listener, err := net.Listen("tcp", ":38808")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Print("listening")

	mux := NewConnMux()
	defer mux.Close()

	go readDevice(dev, mux)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		mux.Add(conn)
	}
}

func readDevice(devicePath string, mux *ConnMux) {
	// serverIP := "localhost:38808"

	f, err := os.Open(devicePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Printf("connected to %s", devicePath)

	b := make([]byte, 24)
	for {
		_, err = f.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		_, err = mux.Write(b)
		if err != nil {
			log.Print(err)
			continue
		}
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
