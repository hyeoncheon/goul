package goul

import "errors"

// constants.
const (
	ErrAdapterReadNotImplemented  = "AdapterReadNotImplemented"
	ErrAdapterWriteNotImplemented = "AdapterWriteNotImplemented"
)

// Adapter is an interface for in/out adapters for pipeline.
type Adapter interface {
	CommonMixin
	Write(in chan Item, message Message) (done chan Item, err error)
	Read(ctrl chan Item, message Message) (out chan Item, err error)
	Close() error
}

//** Base Implementation

// BaseAdapter is a base implementation for the Adapter interface
type BaseAdapter struct {
	BaseCommon
	ID string
}

func (a *BaseAdapter) Read(ctrl chan Item, message Message) (chan Item, error) {
	return nil, errors.New(ErrAdapterReadNotImplemented)
}

func (a *BaseAdapter) Write(in chan Item, message Message) (chan Item, error) {
	return nil, errors.New(ErrAdapterWriteNotImplemented)
}

// Close implements Adapter interface
func (a *BaseAdapter) Close() error {
	return nil
}
