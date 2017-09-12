package pipes

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/hyeoncheon/goul"
)

//** sample implementation of goul.PacketPipe: packet printer -------

// PacketPrinter is a sample pipe print out packet structure to standard out.
// Acceptable data type: pcap formatted packet or raw packet
type PacketPrinter struct{}

// Pipe implements goul.PacketPipe interface
func (p *PacketPrinter) Pipe(in, out chan goul.Item) {
	defer close(out)

	fmt.Println("PacketPrinter ready...")
	for item := range in {
		if p, ok := item.(gopacket.Packet); ok {
			fmt.Println(p)
		} else {
			p := gopacket.NewPacket(item.Data(), layers.LayerTypeEthernet, gopacket.Default)
			if packet, ok := p.(gopacket.Packet); ok {
				fmt.Println(packet)
				out <- packet
				continue
			}
		}
		out <- item
	}
	fmt.Println("PacketPrinter finished.")
}

// Reverse implements goul.PacketPipe interface
func (p *PacketPrinter) Reverse(in, out chan goul.Item) {
	p.Pipe(in, out)
}

//** sample implementation of goul.PacketPipe: data counter ---------

// DataCounter is a sample pipe just count the packets passed through.
// Acceptable data type: ANY
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

// NullWriter is a sample writer pipe but not write anywhere in fact.
// Acceptable data type: ANY
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

// SetLogger sets logger for the goul instance.
func (d *NullWriter) SetLogger(l goul.Logger) error {
	return nil
}
