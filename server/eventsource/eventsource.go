package eventsource

import "net/http"

type EventSource struct {
	w http.ResponseWriter
}

func New(w http.ResponseWriter) *EventSource {
	w.Header().Set("Content-Type", "text/event-stream")
	return &EventSource{w: w}
}

func (s *EventSource) Send(e *Event) error {
	b, err := e.MarshalText()
	if err != nil {
		return err
	}

	_, err = s.w.Write(b)
	if err != nil {
		return err
	}

	_, err = s.w.Write([]byte("\n\n"))
	if err != nil {
		return err
	}

	if f, ok := s.w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}
