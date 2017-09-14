package pipes

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/hyeoncheon/goul"
)

// CompressGZip is a sample processing pipe that compress the packet with gzip.
type CompressGZip struct {
	inLoop bool
	err    error
}

// InLoop implements goul.PacketPipe interface
func (c *CompressGZip) InLoop() bool {
	return c.inLoop
}

// GetError implements goul.PacketPipe interface
func (c *CompressGZip) GetError() error {
	return c.err
}

// Pipe implements goul.PacketPipe interface
func (c *CompressGZip) Pipe(in, out chan goul.Item) {
	defer close(out)

	var count, totOrig, totComp int64
	var b bytes.Buffer
	c.inLoop = true
	fmt.Println("CompressGZip#Pipe ready...")
	for item := range in {
		if item.String() == "application/gzip" {
			fmt.Println("item is not a gzip compressed file: ", item)
			out <- item
			continue
		}
		b.Truncate(0)

		w := gzip.NewWriter(&b)
		w.Write(item.Data())
		w.Flush()
		w.Close()

		sizeOrig := len(item.Data())
		sizeComp := len(b.Bytes())
		fmt.Printf("gzip com size: %v/%v=%.2f\n", sizeComp, sizeOrig, float64(sizeComp)/float64(sizeOrig)*100.0)

		out <- &goul.ItemGeneric{Meta: "application/gzip", DATA: b.Bytes()}

		totOrig += int64(sizeOrig)
		totComp += int64(sizeComp)
		count++
	}
	c.err = errors.New(goul.ErrPipeInputClosed)
	fmt.Printf("CompressGZip#Pipe: total %v packets, %v bytes, %.1f%%\n", count, totOrig, float64(totComp)/float64(totOrig)*100.0)
	c.inLoop = false
}

// Reverse implements goul.PacketPipe interface
func (c *CompressGZip) Reverse(in, out chan goul.Item) {
	defer close(out)

	var count, totOrig, totComp int64
	var b bytes.Buffer
	c.inLoop = true
	fmt.Println("CompressGZip#Reverse ready...")
	for item := range in {
		if item.String() != "application/gzip" {
			fmt.Println("item is not a gzip compressed file: ", item)
			out <- item
			continue
		}

		b.Truncate(0)
		b.Write(item.Data())

		r, err := gzip.NewReader(&b)
		if err != nil {
			fmt.Println("could not create gzip reader:", err)
			c.err = errors.New("CouldNotCreateNewReader")
			continue
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println("ioutil error:", err)
		}
		r.Close()

		sizeOrig := len(item.Data())
		sizeComp := len(buf)
		fmt.Printf("gzip dec size: %v/%v\n", sizeOrig, sizeComp)

		// TODO need to check the type of the buf but... do I deprecate it?
		out <- gopacket.NewPacket(buf, layers.LayerTypeEthernet, gopacket.Default)

		totOrig += int64(sizeOrig)
		totComp += int64(sizeComp)
		count++
	}
	c.err = errors.New(goul.ErrPipeInputClosed)
	fmt.Printf("CompressGZip#Reverse: total %v packets, %v bytes, %.1f%%\n", count, totOrig, float64(totComp)/float64(totOrig)*100.0)
	c.inLoop = false
}
