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
	Write(in chan Item, message Message) (out chan Item, err error)
	Read(in chan Item, message Message) (out chan Item, err error)
}

//** Base Implementation

// BaseAdapter is a base implementation for the Adapter interface
type BaseAdapter struct {
	BaseCommon
	ID string
}

func (a *BaseAdapter) Read(in chan Item, message Message) (chan Item, error) {
	return nil, errors.New(ErrAdapterReadNotImplemented)
}

func (a *BaseAdapter) Write(in chan Item, message Message) (chan Item, error) {
	return nil, errors.New(ErrAdapterWriteNotImplemented)
}
