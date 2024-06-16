package bluetooth

import (
	"github.com/abibby/salusa/set"
	"github.com/godbus/dbus/v5"
)

func (bt *BT) addMatchSignal(options ...dbus.MatchOption) (func() error, error) {
	bt.matchID++
	id := bt.matchID

	toAdd := []dbus.MatchOption{}
	for _, opt := range options {
		m, ok := bt.matches[opt]
		if !ok {
			m = set.New[int]()
			toAdd = append(toAdd, opt)
		}
		m.Add(id)
	}

	err := bt.conn.AddMatchSignal(toAdd...)
	if err != nil {
		return nil, err
	}

	return func() error {
		toRemove := []dbus.MatchOption{}
		for _, opt := range options {
			m, ok := bt.matches[opt]
			if !ok {
				continue
			}
			m.Delete(id)
			if m.Len() == 0 {
				toRemove = append(toRemove, opt)
			}
		}

		return bt.conn.RemoveMatchSignal(toRemove...)
	}, nil
}
