package goul

//** Mixin methods for adapters and pipes.

// CommonMixin is an interface for common functions
type CommonMixin interface {
	SetLogger(logger Logger) error
	GetLogger() Logger
	SetError(err error)
	GetError() error
}

// BaseCommon is a base implementation for the CommonMixin
type BaseCommon struct {
	logger Logger
	err    error
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

// SetError implements CommonMixin
func (c *BaseCommon) SetError(err error) {
	c.err = err
	return
}

// GetError implements CommonMixin
func (c *BaseCommon) GetError() error {
	return c.err
}