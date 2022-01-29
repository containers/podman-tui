package utils

// OutputStreamChannel implements io.WriterCloser with channel.
type OutputStreamChannel struct {
	channel chan string
}

// NewStreamChannel returns a new io.WriterCloser with channel.
// This will attach to container session to display results on terminal dialog primitive.
func NewStreamChannel(size int) *OutputStreamChannel {
	stream := OutputStreamChannel{
		channel: make(chan string, size),
	}
	return &stream
}

// Channel returns outputStreamChannel channel address
func (osc *OutputStreamChannel) Channel() *chan string {
	return &osc.channel
}

// Write writes to channel.
func (osc *OutputStreamChannel) Write(p []byte) (int, error) {
	osc.channel <- string(p)
	return len(p), nil
}

// Close stream closer
func (osc *OutputStreamChannel) Close() error {
	return nil
}
