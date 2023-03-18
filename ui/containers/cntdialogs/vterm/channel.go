package vterm

import (
	"errors"
	"io"
	"sync"
)

// Writer is an io.Writer that proxies Write() calls to a channel
// The []byte buffer of the Write() is queued on the channel as one message.
type Writer interface {
	io.Writer
	Chan() <-chan []byte
}

type writer struct {
	ch  chan []byte
	mux sync.Mutex
}

// NewWriter initializes a new channel writer
func NewWriter(c chan []byte) Writer {
	return &writer{
		ch: c,
	}
}

// Chan returns the R/O channel behind Writer
func (w *writer) Chan() <-chan []byte {
	return w.ch
}

// Write method for Writer
func (w *writer) Write(b []byte) (bLen int, err error) {
	if w == nil {
		return 0, errors.New("use channel.NewWriter() to initialize a Writer")
	}

	w.mux.Lock()
	defer w.mux.Unlock()

	if w.ch == nil {
		return 0, errors.New("the channel is closed for Write")
	}

	buf := make([]byte, len(b))
	copy(buf, b)
	w.ch <- buf

	return len(b), nil
}
