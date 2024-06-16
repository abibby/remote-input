package eventsource

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

type Event struct {
	// A string identifying the type of event described. If this is specified,
	// an event will be dispatched on the browser to the listener for the
	// specified event name; the website source code should use
	// addEventListener() to listen for named events. The onmessage handler is
	// called if no event name is specified for a message.
	Event string

	// The data field for the message. When the EventSource receives multiple
	// consecutive lines that begin with data:, it concatenates them, inserting
	// a newline character between each one. Trailing newlines are removed.
	Data string

	// The event ID to set the EventSource object's last event ID value.
	ID string

	// The reconnection time. If the connection to the server is lost, the
	// browser will wait for the specified time before attempting to reconnect.
	// This must be an integer, specifying the reconnection time in
	// milliseconds. If a non-integer value is specified, the field is ignored.
	Retry int
}

var _ encoding.TextMarshaler = (*Event)(nil)

// MarshalText implements encoding.TextMarshaler.
func (e *Event) MarshalText() (text []byte, err error) {
	b := []byte{}
	b = field(b, "event", e.Event)
	b = field(b, "data", e.Data)
	b = field(b, "id", e.ID)
	b = field(b, "retry", e.Retry)
	return b, nil
}

func field(b []byte, name string, value any) []byte {
	if reflect.ValueOf(value).IsZero() {
		return b
	}
	for _, line := range strings.Split(fmt.Sprint(value), "\n") {
		b = append(b, []byte(name+": "+line+"\n")...)
	}
	return b
}
