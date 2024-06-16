package bluetooth

import (
	"fmt"

	"github.com/abibby/salusa/set"
	"github.com/godbus/dbus/v5"
)

type BT struct {
	conn    *dbus.Conn
	matchID int
	matches map[dbus.MatchOption]set.Set[int]
}

func New() (*BT, error) {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system bus: %w", err)
	}

	return &BT{
		conn:    conn,
		matches: map[dbus.MatchOption]set.Set[int]{},
	}, nil
}

func (bt *BT) Close() error {
	return bt.conn.Close()
}
