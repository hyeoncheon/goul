package pipes

import (
	"fmt"

	"github.com/google/gopacket"

	"github.com/hyeoncheon/goul"
)

//** sample filters and writers -------------------------------------

// PacketPrinter is a sample pipe print out packet structure to standard out.
type PacketPrinter struct{}

// Pipe implements goul.PacketPipe interface
func (p *PacketPrinter) Pipe(in, out chan goul.Item) {
	defer close(out)

	fmt.Println("PacketPrinter ready...")
	for item := range in {
		if p, ok := item.(gopacket.Packet); ok {
			fmt.Println(p)
		}
		out <- item
	}
	fmt.Println("PacketPrinter finished.")
}

// Reverse implements goul.PacketPipe interface
func (p *PacketPrinter) Reverse(in, out chan goul.Item) {
	p.Pipe(in, out)
}

// DataCounter is a sample pipe just count the packets passed through.
type DataCounter struct{}

// Pipe implements goul.PacketPipe interface
func (c *DataCounter) Pipe(in, out chan goul.Item) {
	defer close(out)

	var count int64
	fmt.Println("DataCounter ready...")
	for item := range in {
		fmt.Println("DataCounter:size: ", len(item.Data()))
		count++
		out <- item
	}
	fmt.Printf("DataCounter counts total %v packets. counter finished.\n", count)
}

// Reverse implements goul.PacketPipe interface
func (c *DataCounter) Reverse(in, out chan goul.Item) {
	c.Pipe(in, out)
}

//** sample implementation of goul.Writer for debugging. ------------

// NullWriter is a sample pipe of simple packet counter
type NullWriter struct{}

// Writer implements Writer interface
func (d *NullWriter) Writer(in chan goul.Item) {
	fmt.Println("NullWriter#Writer ready...")

	var count int64
	for range in {
		count++
	}
	fmt.Printf("NullWriter#Writer counts total %v packets. counter finished.\n", count)
}
