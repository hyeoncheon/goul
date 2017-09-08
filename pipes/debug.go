package pipes

import (
	"fmt"

	"github.com/google/gopacket"

	"github.com/hyeoncheon/goul"
)

//** sample filters and writers -------------------------------------

// PacketPrinter is a sample pipe of simple standard out.
func PacketPrinter(in, out chan goul.Item) {
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

// DataCounter is a sample pipe of simple standard out.
func DataCounter(in, out chan goul.Item) {
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

// DataWriter is a sample pipe of simple packet counter
func DataWriter(in chan goul.Item) {
	fmt.Println("DataWriter ready...")

	var count int64
	for _ = range in {
		count++
	}
	fmt.Printf("DataWriter counts total %v packets. counter finished.\n", count)
}
