package services

import (
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"syscall"
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

var _ io.WriteCloser = (*ConnMux)(nil)

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
			if err != io.EOF || isConnReset(err) {
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

func isConnReset(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			if sysErr.Err == syscall.ECONNRESET {
				return true
			}
		}
	}
	return false
}
