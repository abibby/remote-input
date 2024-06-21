package eventsource

import (
	"io"
	"log"
	"net/http"
)

type EventSource struct {
	events       chan *Event
	errorHandler func(err error)
}

func New(events chan *Event) *EventSource {
	return &EventSource{
		events: events,
		errorHandler: func(err error) {
			log.Print(err)
		},
	}
}

func (s *EventSource) SetErrorHandler(errorHandler func(error)) *EventSource {
	s.errorHandler = errorHandler
	return s
}
func (s *EventSource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	for e := range s.events {
		err := send(w, e)
		if err != nil {
			s.errorHandler(err)
		}
	}
}

func send(w io.Writer, e *Event) error {
	b, err := e.MarshalText()
	if err != nil {
		return err
	}

	_, err = w.Write(append(b, '\n', '\n'))
	if err != nil {
		return err
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}
