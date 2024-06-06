package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/abibby/remote-input/common"
	"github.com/abibby/remote-input/windows"
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
		if err == io.EOF {
			toRemove = append(toRemove, c)
		} else if err != nil {
			errs = append(errs, err)
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

		// fmt.Printf("%v", len(b))

		// e := &common.InputEvent{}

		// err = e.UnmarshalBinary(b)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		_, err = mux.Write(b)
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

var keyMap = map[common.KeyCode]windows.VirtualKey{
	common.KEY_BACKSPACE: windows.VK_BACK,
	common.KEY_ENTER:     windows.VK_RETURN,
	common.KEY_TAB:       windows.VK_TAB,
	common.KEY_CAPSLOCK:  windows.VK_CAPITAL,

	common.KEY_LEFTCTRL:   windows.VK_CONTROL,
	common.KEY_RIGHTCTRL:  windows.VK_CONTROL,
	common.KEY_LEFTSHIFT:  windows.VK_SHIFT,
	common.KEY_RIGHTSHIFT: windows.VK_SHIFT,

	common.KEY_A: windows.VK_A,
	common.KEY_B: windows.VK_B,
	common.KEY_C: windows.VK_C,
	common.KEY_D: windows.VK_D,
	common.KEY_E: windows.VK_E,
	common.KEY_F: windows.VK_F,
	common.KEY_G: windows.VK_G,
	common.KEY_H: windows.VK_H,
	common.KEY_I: windows.VK_I,
	common.KEY_J: windows.VK_J,
	common.KEY_K: windows.VK_K,
	common.KEY_L: windows.VK_L,
	common.KEY_M: windows.VK_M,
	common.KEY_N: windows.VK_N,
	common.KEY_O: windows.VK_O,
	common.KEY_P: windows.VK_P,
	common.KEY_Q: windows.VK_Q,
	common.KEY_R: windows.VK_R,
	common.KEY_S: windows.VK_S,
	common.KEY_T: windows.VK_T,
	common.KEY_U: windows.VK_U,
	common.KEY_V: windows.VK_V,
	common.KEY_W: windows.VK_W,
	common.KEY_X: windows.VK_X,
	common.KEY_Y: windows.VK_Y,
	common.KEY_Z: windows.VK_Z,

	common.KEY_1: windows.VK_1,
	common.KEY_2: windows.VK_2,
	common.KEY_3: windows.VK_3,
	common.KEY_4: windows.VK_4,
	common.KEY_5: windows.VK_5,
	common.KEY_6: windows.VK_6,
	common.KEY_7: windows.VK_7,
	common.KEY_8: windows.VK_8,
	common.KEY_9: windows.VK_9,
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
