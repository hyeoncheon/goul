package goul

import "errors"

// Router is a centural unit of the data processing pipeline for both
// client and server. By configuring proper reader and writer handlers,
// the program can act as capturing sender receiving client.

// error messages of router subsystem.
const (
	ErrRouterPipelineNotSupported = "RouterPipelineNotSupported"
	ErrRouterNoReaderOrWriter     = "RouterNoReaderOrWriter"
)

// Router is an interface for the heart of the goul processing.
//
type Router interface {
	// Run prepares all necessary things and start reader, writer and all
	// registered processors.
	// BaseRouter provides simple Run function with reader and writer but
	// You should implement this function for your specific router.
	Run() (control, done chan Item, err error)
	// AddPipe adds given data processor to the end of the processor list.
	// BaseRouter does not support pipeline and you should implement your
	// own AddPipe function for your router.
	AddPipe(pipe Pipe) error

	// SetReader sets given reader as a data source of the processing router.
	SetReader(reader Adapter) error
	// SetWriter sets given writer as a data sink of the processing router.
	SetWriter(writer Adapter) error
	// SetLogger sets given logger as a logger for the router.
	SetLogger(logger Logger) error
	// GetPipes returns a slice of pipe handlers.
	GetPipes() []Pipe
	getReader() Adapter
	getWriter() Adapter
	getLogger() Logger
}

// BaseRouter is a simple/sample router that just support in/out adapters.
// You can use it as a base of your own router implementation so you can
// just implement and override interesting parts of your implementation.
// In this case, you can use following statement:
//
//	router := &MyRouter{Router: &BaseRouter{}}
//
// Please refer to goul.Pipeline and related test cases for more detail.
// (router_pipeline.go)
type BaseRouter struct {
	err    error
	logger Logger
	reader Adapter
	writer Adapter
}

// Run implements Router.
// This implementation just invoke reader and writer without pipes.
func (r *BaseRouter) Run() (chan Item, chan Item, error) {
	if r.getReader() == nil || r.getWriter() == nil {
		return nil, nil, errors.New(ErrRouterNoReaderOrWriter)
	}

	cntl := make(chan Item)
	ch, err := r.getReader().Read(cntl, nil)
	if err != nil {
		return nil, nil, err
	}
	tx, err := r.getWriter().Write(ch, nil)
	if err != nil {
		return nil, nil, err
	}

	Log(r.getLogger(), "ROUTER", "--------------- started -----------------")
	return cntl, tx, nil
}

// AddPipe implements Router: but BaseRouter does not support pipelining.
func (r *BaseRouter) AddPipe(pipe Pipe) error {
	return errors.New(ErrRouterPipelineNotSupported)
}

// GetPipes implements Router: placeholder
func (r *BaseRouter) GetPipes() []Pipe {
	return []Pipe{}
}

// SetReader implements Router: working
func (r *BaseRouter) SetReader(reader Adapter) error {
	reader.SetLogger(r.logger)
	r.reader = reader
	return nil
}

func (r *BaseRouter) getReader() Adapter {
	return r.reader
}

// SetWriter implements Router: working
func (r *BaseRouter) SetWriter(writer Adapter) error {
	writer.SetLogger(r.logger)
	r.writer = writer
	return nil
}

func (r *BaseRouter) getWriter() Adapter {
	return r.writer
}

// SetLogger implements Router: working
func (r *BaseRouter) SetLogger(logger Logger) error {
	r.logger = logger
	return nil
}

func (r *BaseRouter) getLogger() Logger {
	return r.logger
}

// Messages is a map for system messages.
var Messages = map[string]*ItemGeneric{
	"closed": &ItemGeneric{
		Meta: "message",
		DATA: []byte("channel closed. done"),
	},
}
