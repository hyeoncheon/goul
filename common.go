package goul

// ChannelSize is a default size of pipeline channels.
const ChannelSize = 10

//** common utilities for pipe handlers

// Message ...
type Message interface {
}

// PipeFunction ...
type PipeFunction func(in, out chan Item, message Message)

// Launch is a helper function that make goroutine execution simple.
func Launch(fn PipeFunction, in chan Item, message Message) (chan Item, error) {
	out := make(chan Item, ChannelSize)
	go fn(in, out, message)
	return out, nil
}
