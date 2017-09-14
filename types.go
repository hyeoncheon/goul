package goul

//** types for goul, items ------------------------------------------

// Item is an interface for passing data between pipes.
type Item interface {
	String() string
	Data() []byte
}

// ItemGeneric is a structure for the generic byte slice data.
type ItemGeneric struct {
	Meta string
	DATA []byte
}

// String implements goul.Item
func (c *ItemGeneric) String() string {
	return c.Meta
}

// Data implements goul.Item
func (c *ItemGeneric) Data() []byte {
	return c.DATA
}

//** types for goul, pipeline functions -----------------------------

// PacketPipe is an interface for the packet/data processing pipe
type PacketPipe interface {
	InLoop() bool
	GetError() error
	Pipe(in, out chan Item)
	Reverse(in, out chan Item)
}

// Reader is an interface for the reader pipe
type Reader interface {
	SetLogger(logger Logger) error
	InLoop() bool
	GetError() error
	Reader(in chan int, out chan Item)
}

// Writer is an interface for the writer pipe
type Writer interface {
	SetLogger(logger Logger) error
	InLoop() bool
	GetError() error
	Writer(in chan Item)
}
