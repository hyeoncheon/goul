package goul

//** Mixin methods for adapters and pipes.

// CommonMixin is an interface for common functions
type CommonMixin interface {
	SetLogger(logger Logger) error
	GetLogger() Logger
}

// BaseCommon is a base implementation for the CommonMixin
type BaseCommon struct {
	logger Logger
}

// SetLogger implements CommonMixin
func (c *BaseCommon) SetLogger(logger Logger) error {
	c.logger = logger
	return nil
}

// GetLogger implements CommonMixin
func (c *BaseCommon) GetLogger() Logger {
	return c.logger
}

//** common utilities for pipe handlers

// Message ...
type Message interface {
}

// PipeFunction ...
type PipeFunction func(in, out chan Item, message Message)

// Launch is a helper function that make goroutine execution simple.
func Launch(fn PipeFunction, in chan Item, message Message) (chan Item, error) {
	out := make(chan Item)
	go fn(in, out, message)
	return out, nil
}
