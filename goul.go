package goul

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// constants
const (
	SigINT = 2

	ErrCannotSetFilter  = "cannot set filter"
	ErrNoReaderOrWriter = "no reader or writer"
)

//** types for goul, items ------------------------------------------

// Item is an interface for passing data between pipes.
type Item interface {
	String() string
	Data() []byte
}

//** type for goul, pipeline functions ------------------------------

// PacketPipe is an interface for the packet/data processing pipe
type PacketPipe interface {
	Pipe(in, out chan Item)
	Reverse(in, out chan Item)
}

// Reader is an interface for the reader pipe
type Reader interface {
	Reader(in chan int, out chan Item)
}

// Writer is an interface for the writer pipe
type Writer interface {
	Writer(in chan Item)
}

// Goul is...
type Goul struct {
	device      string
	snaplen     int
	promiscuous bool
	timeout     time.Duration
	inactive    *pcap.InactiveHandle
	handle      *pcap.Handle
	filter      string

	isReceiver bool
	err        error
	logger     Logger

	//	reader func(in chan int, out chan gopacket.Packet)
	cmdChannel chan int
	bufferSize int
	reader     Reader
	writer     Writer
	pipes      []PacketPipe
}

// New is...
func New(dev string, isReceiver bool, cmd chan int) (*Goul, error) {
	g := &Goul{
		device:      dev,
		snaplen:     1600,
		promiscuous: false,
		timeout:     1,
		isReceiver:  isReceiver,
		cmdChannel:  cmd,
		bufferSize:  1,
		filter:      "ip",
	}
	g.inactive, g.err = pcap.NewInactiveHandle(g.device)
	if g.isReceiver {
		g.writer = g
	} else {
		g.reader = g
	}
	return g, g.err
}

// Close is...
func (g *Goul) Close() {
	if g.handle != nil {
		g.handle.Close()
	}
	if g.inactive != nil {
		g.inactive.CleanUp()
	}
}

// Run checks the pipeline condition and execute all pipes including reader,
// then run writer.
func (g *Goul) Run() error {
	if g.reader == nil || g.writer == nil {
		return errors.New(ErrNoReaderOrWriter)
	}
	ch := g.ExecReader()
	for _, pipe := range g.pipes {
		ch = g.ExecPipe(pipe.Pipe, ch)
	}
	g.writer.Writer(ch)
	return nil
}

// ExecReader creates new channel for output channel of the reader and execute
// the reader as goroutine. It returns created output channel so next pipe can
// be connected to this pipeline.
func (g *Goul) ExecReader() chan Item {
	ch := make(chan Item, g.bufferSize)
	go g.reader.Reader(g.cmdChannel, ch)
	return ch
}

// ExecPipe creates new channel for output channel of the pipe and execute
// the pipe as goroutine. It returns created output channel so next pipe can
// be connected to this pipeline.
func (g *Goul) ExecPipe(fn func(in, out chan Item), in chan Item) chan Item {
	ch := make(chan Item, g.bufferSize)
	go fn(in, ch)
	return ch
}

// SetOptions is...
func (g *Goul) SetOptions(promiscuous bool, len int, timeout time.Duration) error {
	if g.err = g.inactive.SetTimeout(timeout * time.Second); g.err != nil {
		g.log("ERROR:", g.err)
	} else if g.err = g.inactive.SetSnapLen(len); g.err != nil {
		g.log("ERROR:", g.err)
	} else if g.err = g.inactive.SetPromisc(promiscuous); g.err != nil {
		g.log("ERROR:", g.err)
	}
	g.promiscuous = promiscuous
	g.snaplen = len
	g.timeout = timeout
	return g.err
}

// SetFilter is... NOT IMPLEMENT YET
func (g *Goul) SetFilter(filter string) error {
	g.filter = filter
	return nil
}

// SetReader is...
func (g *Goul) SetReader(r Reader) error {
	g.reader = r
	return nil
}

// SetWriter is...
func (g *Goul) SetWriter(w Writer) error {
	g.writer = w
	return nil
}

// AddPipe is...
func (g *Goul) AddPipe(pipe PacketPipe) error {
	g.pipes = append(g.pipes, pipe)
	return nil
}

// SetLogger sets logger for the goul instance.
func (g *Goul) SetLogger(l Logger) error {
	g.logger = l
	return nil
}

func (g *Goul) log(args ...interface{}) {
	if g.logger != nil {
		g.logger.Debug(args...)
	}
}

func (g *Goul) logf(format string, args ...interface{}) {
	if g.logger != nil {
		g.logger.Debugf(format, args...)
	}
}

//** tap methods... -------------------------------------------------

// Reader implements Reader interface
func (g *Goul) Reader(cmd chan int, out chan Item) {
	g.Capture(cmd, out)
}

// Capture is a device packet reader which used as reader when capturer mode.
func (g *Goul) Capture(cmd chan int, out chan Item) {
	defer close(out)

	if g.handle == nil && g.inactive != nil {
		g.log("capture: not activated. activating...")
		g.handle, g.err = g.inactive.Activate()
		if g.err != nil {
			g.log("cannot activate handler:", g.err)
			return
		}
		defer g.handle.Close()
	}
	if g.err = g.handle.SetBPFFilter(g.filter); g.err != nil {
		g.logf("cannot set filter: %v", g.err)
		return
	}

	fmt.Println("capture ready...")
	packetSource := gopacket.NewPacketSource(g.handle, g.handle.LinkType())
	for {
		select {
		case command := <-cmd:
			switch command {
			case SigINT:
				fmt.Println("interrupt! shutdown capture.")
				return
			}
		case packet := <-packetSource.Packets():
			out <- packet
		default: // for unblock from the channel
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// Writer implements Writer interface
func (g *Goul) Writer(in chan Item) {
	g.Inject(in)
}

// Inject is a device packet writer which used as writer when injector mode.
func (g *Goul) Inject(in chan Item) {
	if g.handle == nil && g.inactive != nil {
		g.log("inject: not activated. activating...")
		g.handle, g.err = g.inactive.Activate()
		if g.err != nil {
			g.log("cannot activate handler:", g.err)
			return
		}
		defer g.handle.Close()
	}

	fmt.Println("Inject ready...")
	for item := range in {
		if p, ok := item.(gopacket.Packet); ok {
			g.handle.WritePacketData(p.Data())
		}
	}
}
