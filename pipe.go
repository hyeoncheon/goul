package goul

import "errors"

// constants.
const (
	ModeConverter = true
	ModeReverter  = false

	ErrPipeConvertNotImplemented = "PipeConvertNotImplemented"
	ErrPipeRevertNotImplemented  = "PipeRevertNotImplemented"
)

// Pipe is an interface for pipeline handlers.
type Pipe interface {
	CommonMixin
	Convert(in chan Item, message Message) (out chan Item, err error)
	Revert(in chan Item, message Message) (out chan Item, err error)
	IsConverter() bool
}

//** Base Implementation

// BasePipe is a base implementation for the Pipe interface
type BasePipe struct {
	BaseCommon
	Mode bool
}

// Convert implements interface Converter
func (p *BasePipe) Convert(in chan Item, message Message) (chan Item, error) {
	return nil, errors.New(ErrPipeConvertNotImplemented)
}

// Revert implements interface Revert
func (p *BasePipe) Revert(in chan Item, message Message) (chan Item, error) {
	return nil, errors.New(ErrPipeRevertNotImplemented)
}

// IsConverter implements Pipe
func (p *BasePipe) IsConverter() bool {
	return p.Mode
}
