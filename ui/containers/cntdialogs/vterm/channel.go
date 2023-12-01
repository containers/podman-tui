package vterm

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrChannelInit   = errors.New("use channel.NewWriter() to initialize a Writer")
	ErrChannelClosed = errors.New("the channel is closed for Write")
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

// NewWriter initializes a new channel writer.
func NewWriter(c chan []byte) Writer { //nolint:ireturn
	return &writer{
		ch: c,
	}
}

// Chan returns the R/O channel behind Writer.
func (w *writer) Chan() <-chan []byte {
	return w.ch
}

// Write method for Writer.
func (w *writer) Write(b []byte) (int, error) {
	if w == nil {
		return 0, ErrChannelInit
	}

	w.mux.Lock()
	defer w.mux.Unlock()

	if w.ch == nil {
		return 0, ErrChannelClosed
	}

	buf := make([]byte, len(b))
	copy(buf, b)
	w.ch <- buf

	return len(b), nil
}
