package goul

import (
	"errors"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// constants
const (
	NAME         = "Goul"
	ComInterrupt = 2

	ItemTypeUnknown   = "unknown"
	ItemTypeRawPacket = "rawpacket"

	ErrPipeInterrupted        = "pipe interrupted"
	ErrPipeInputClosed        = "input channel closed"
	ErrPipeOutputClosed       = "output channel closed"
	ErrNetworkReadHeader      = "could not read header from network"
	ErrNetworkReadChunk       = "could not read chunk from network"
	ErrNetworkConnectionReset = "connection reset"
	ErrCannotSetFilter        = "could not set filter"
	ErrNoReaderOrWriter       = "no reader or writer"
)

// Goul is base structure of this masul goul (magic mirror).
type Goul struct {
	device      string
	snaplen     int
	promiscuous bool
	timeout     time.Duration
	inactive    *pcap.InactiveHandle
	handle      *pcap.Handle
	filter      string

	isTest     bool
	isReceiver bool
	err        error
	inLoop     bool
	logger     Logger

	//	reader func(in chan int, out chan gopacket.Packet)
	cmdChannel chan int
	bufferSize int
	reader     Reader
	writer     Writer
	pipes      []PacketPipe
}

// New create new Goul instance and setup inactive handler of pcap.
func New(dev string, recv bool, cmd chan int, test bool) (*Goul, error) {
	g := &Goul{
		device:      dev,
		snaplen:     1600,
		promiscuous: false,
		timeout:     1,
		isTest:      test,
		isReceiver:  recv,
		cmdChannel:  cmd,
		bufferSize:  1,
		filter:      "ip",
	}
	g.inactive, g.err = pcap.NewInactiveHandle(g.device)
	if g.isReceiver {
		g.log("preparing as receiver mode...")
		g.writer = g
	} else {
		g.log("preparing as sender mode...")
		g.reader = g
	}
	return g, g.err
}

// Close clean up resources.
func (g *Goul) Close() {
	g.log("cleanup...")
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
		g.err = errors.New(ErrNoReaderOrWriter)
		return g.err
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

// SetOptions sets capture options to inactive handler.
func (g *Goul) SetOptions(promiscuous bool, len int, timeout time.Duration) error {
	if g.err = g.inactive.SetTimeout(timeout * time.Second); g.err != nil {
		g.log("ERROR: ", g.err)
	} else if g.err = g.inactive.SetSnapLen(len); g.err != nil {
		g.log("ERROR: ", g.err)
	} else if g.err = g.inactive.SetPromisc(promiscuous); g.err != nil {
		g.log("ERROR: ", g.err)
	}
	g.promiscuous = promiscuous
	g.snaplen = len
	g.timeout = timeout
	return g.err
}

// SetFilter sets filter string which is applied while capturing.
func (g *Goul) SetFilter(filter string) error {
	g.filter = filter
	return nil
}

// SetReader sets Reader object (otherwise, Goul will be used)
func (g *Goul) SetReader(r Reader) error {
	g.reader = r
	return nil
}

// SetWriter sets Writer object (otherwise, Goul will be used)
func (g *Goul) SetWriter(w Writer) error {
	g.writer = w
	return nil
}

// AddPipe adds a pipeline object into pipeline stack.
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
		args = append([]interface{}{NAME + ": "}, args...)
		g.logger.Debug(args...)
	}
}

func (g *Goul) logf(format string, args ...interface{}) {
	if g.logger != nil {
		g.logger.Debugf(NAME+": "+format, args...)
	}
}

//** tap methods... -------------------------------------------------

// InLoop implements goul.PacketPipe interface
func (g *Goul) InLoop() bool {
	return g.inLoop
}

// GetError implements Writer/Reader interface
func (g *Goul) GetError() error {
	return g.err
}

// Reader implements Reader interface (wrapper)
func (g *Goul) Reader(cmd chan int, out chan Item) {
	err := g.Capture(cmd, out)
	if err != nil {
		g.logger.Error("could not start capture: ", err)
	}
}

// Capture is a device packet reader which used as reader when capturer mode.
func (g *Goul) Capture(cmd chan int, out chan Item) error {
	defer close(out)

	if g.handle == nil && g.inactive != nil {
		g.log("capture: not activated. activating...")
		g.handle, g.err = g.inactive.Activate()
		if g.err != nil {
			g.logger.Error("cannot activate handler: ", g.err)
			return g.err
		}
		defer g.handle.Close()
	}
	if g.err = g.handle.SetBPFFilter(g.filter); g.err != nil {
		g.log("cannot set filter: ", g.err)
		return g.err
	}

	packetSource := gopacket.NewPacketSource(g.handle, g.handle.LinkType())
	g.inLoop = true
	g.log("capturing started...")
	for {
		select {
		case command := <-cmd:
			switch command {
			case ComInterrupt:
				g.log("shutting down capturing...")
				g.err = errors.New(ErrPipeInterrupted)
				g.inLoop = false
				return nil
			}
		case packet := <-packetSource.Packets():
			out <- packet
		default: // for unblock from the channel
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Writer implements Writer interface
func (g *Goul) Writer(in chan Item) {
	if g.isTest {
		g.log("dummy writer ready...")

		var count int64
		g.inLoop = true
		for range in {
			count++
		}
		g.inLoop = false
		g.logf("dummy writer counts total %v packets. exit.", count)
	} else {
		// it's real writer. yes.
		g.Inject(in)
	}
}

// Inject is a device packet writer which used as writer when injector mode.
func (g *Goul) Inject(in chan Item) {
	if g.handle == nil && g.inactive != nil {
		g.log("handle not activated. activating...")
		g.handle, g.err = g.inactive.Activate()
		if g.err != nil {
			g.log("cannot activate handler: ", g.err)
			return
		}
		defer g.handle.Close()
	}

	g.inLoop = true
	g.log("injection started...")
	for item := range in {
		if p, ok := item.(gopacket.Packet); ok {
			g.handle.WritePacketData(p.Data())
		}
	}
	g.inLoop = false
}
