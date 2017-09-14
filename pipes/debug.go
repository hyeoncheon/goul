package pipes

import (
	"errors"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/hyeoncheon/goul"
)

//** sample implementation of goul.PacketPipe: packet printer -------

// PacketPrinter is a sample pipe print out packet structure to standard out.
// Acceptable data type: pcap formatted packet or raw packet
type PacketPrinter struct {
	inLoop bool
	err    error
}

// InLoop implements goul.PacketPipe interface
func (p *PacketPrinter) InLoop() bool {
	return p.inLoop
}

// GetError implements goul.PacketPipe interface
func (p *PacketPrinter) GetError() error {
	return p.err
}

// Pipe implements goul.PacketPipe interface
func (p *PacketPrinter) Pipe(in, out chan goul.Item) {
	defer close(out)

	p.inLoop = true
	fmt.Println("PacketPrinter ready...")
	for item := range in {
		if p, ok := item.(gopacket.Packet); ok {
			fmt.Println(p)
		} else if item.String() == goul.ItemTypeRawPacket {
			p := gopacket.NewPacket(item.Data(), layers.LayerTypeEthernet, gopacket.Default)
			if packet, ok := p.(gopacket.Packet); ok {
				fmt.Println(packet)
				out <- packet
				continue
			}
		}
		out <- item
	}
	p.err = errors.New(goul.ErrPipeInputClosed)
	fmt.Println("PacketPrinter finished.")
	p.inLoop = false
}

// Reverse implements goul.PacketPipe interface
func (p *PacketPrinter) Reverse(in, out chan goul.Item) {
	p.Pipe(in, out)
}

//** sample implementation of goul.PacketPipe: data counter ---------

// DataCounter is a sample pipe just count the packets passed through.
// Acceptable data type: ANY
type DataCounter struct {
	inLoop bool
	err    error
}

// InLoop implements goul.PacketPipe interface
func (c *DataCounter) InLoop() bool {
	return c.inLoop
}

// GetError implements goul.PacketPipe interface
func (c *DataCounter) GetError() error {
	return c.err
}

// Pipe implements goul.PacketPipe interface
func (c *DataCounter) Pipe(in, out chan goul.Item) {
	defer close(out)

	var count int64
	c.inLoop = true
	fmt.Println("DataCounter ready...")
	for item := range in {
		fmt.Println("DataCounter:size: ", len(item.Data()))
		count++
		out <- item
	}
	c.err = errors.New(goul.ErrPipeInputClosed)
	fmt.Printf("DataCounter counts total %v packets. counter finished.\n", count)
	c.inLoop = false
}

// Reverse implements goul.PacketPipe interface
func (c *DataCounter) Reverse(in, out chan goul.Item) {
	c.Pipe(in, out)
}

//** sample implementation of goul.Writer for debugging. ------------

// NullWriter is a sample writer pipe but not write anywhere in fact.
// Acceptable data type: ANY
type NullWriter struct {
	inLoop bool
	err    error
}

// InLoop implements goul.PacketPipe interface
func (w *NullWriter) InLoop() bool {
	return w.inLoop
}

// GetError implements goul.PacketPipe interface
func (w *NullWriter) GetError() error {
	return w.err
}

// Writer implements Writer interface
func (w *NullWriter) Writer(in chan goul.Item) {
	fmt.Println("NullWriter#Writer ready...")

	var count int64
	w.inLoop = true
	for range in {
		count++
	}
	w.err = errors.New(goul.ErrPipeInputClosed)
	fmt.Printf("NullWriter#Writer counts total %v packets. counter finished.\n", count)
	w.inLoop = false
}

// SetLogger sets logger for the goul instance. (dummy)
func (w *NullWriter) SetLogger(l goul.Logger) error {
	return nil
}
