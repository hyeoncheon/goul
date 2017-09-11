package pipes

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/hyeoncheon/goul"
)

// CompressGZip is a sample processing pipe that compress the packet with gzip.
type CompressGZip struct{}

// Pipe implements goul.PacketPipe interface
func (c *CompressGZip) Pipe(in, out chan goul.Item) {
	defer close(out)

	var count, totOrig, totComp int64
	var b bytes.Buffer
	fmt.Println("CompressGZip#Pipe ready...")
	for item := range in {
		b.Truncate(0)

		w := gzip.NewWriter(&b)
		w.Write(item.Data())
		w.Flush()
		w.Close()

		sizeOrig := len(item.Data())
		sizeComp := len(b.Bytes())
		fmt.Printf("gzip com size: %v/%v=%.2f\n", sizeComp, sizeOrig, float64(sizeComp)/float64(sizeOrig)*100.0)

		out <- &ItemGeneric{data: b.Bytes()}

		totOrig += int64(sizeOrig)
		totComp += int64(sizeComp)
		count++
	}
	fmt.Printf("CompressGZip#Pipe: total %v packets, %v bytes, %.1f%%\n", count, totOrig, float64(totComp)/float64(totOrig)*100.0)
}

// Reverse implements goul.PacketPipe interface
func (c *CompressGZip) Reverse(in, out chan goul.Item) {
	defer close(out)

	var count, totOrig, totComp int64
	var b bytes.Buffer
	fmt.Println("CompressGZip#Reverse ready...")
	for item := range in {
		b.Truncate(0)
		b.Write(item.Data())

		r, err := gzip.NewReader(&b)
		if err != nil {
			fmt.Println("zlib read error", err)
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println("ioutil error", err)
		}
		r.Close()

		sizeOrig := len(item.Data())
		sizeComp := len(buf)
		fmt.Printf("gzip dec size: %v/%v\n", sizeOrig, sizeComp)

		out <- gopacket.NewPacket(buf, layers.LayerTypeEthernet, gopacket.Default)

		totOrig += int64(sizeOrig)
		totComp += int64(sizeComp)
		count++
	}
	fmt.Printf("CompressGZip#Reverse: total %v packets, %v bytes, %.1f%%\n", count, totOrig, float64(totComp)/float64(totOrig)*100.0)
}
