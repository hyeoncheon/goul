package pipes

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/hyeoncheon/goul"
)

// constants...
const (
	CouldNotCreateNewZLibReader = "couldn't create new reader"
)

// CompressZLib is a sample processing pipe that compress the packet with zlib.
type CompressZLib struct {
	goul.Pipe
	ID string
}

// Convert implements interface Pipe/Converter
func (p *CompressZLib) Convert(in chan goul.Item, message goul.Message) (out chan goul.Item, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "CompressZLib#Convert recovered from panic!\n")
			fmt.Fprintf(os.Stderr, "Probably an inheritance problem of pipeline instance.\n")
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = errors.New("panic")
		}
	}()

	if p.ID == "" {
		p.ID = "zlib-convert"
	}
	p.SetError(nil)
	return goul.Launch(p.converter, in, message)
}

// Revert implements interface Pipe/Reverter
func (p *CompressZLib) Revert(in chan goul.Item, message goul.Message) (out chan goul.Item, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "CompressZLib#Revert recovered from panic!\n")
			fmt.Fprintf(os.Stderr, "Probably an inheritance problem of pipeline instance.\n")
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = errors.New("panic")
		}
	}()

	if p.ID == "" {
		p.ID = "zlib-revert"
	}
	p.SetError(nil)
	return goul.Launch(p.reverter, in, message)
}

// converter compresses items from `in` channel and put it into `out` channel.
func (p *CompressZLib) converter(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(p.GetLogger(), p.ID, "exit")

	var count, totOrig, totComp int64
	var b bytes.Buffer
	goul.Log(p.GetLogger(), p.ID, "zlib compressor in looping...")
	for item := range in {
		if item.String() == "application/zlib" {
			goul.Log(p.GetLogger(), p.ID, "item is not a zlib compressed file: %v", item)
			out <- item
			continue
		}
		b.Truncate(0)

		w := zlib.NewWriter(&b)
		w.Write(item.Data())
		w.Flush()
		w.Close()

		sizeOrig := len(item.Data())
		sizeComp := len(b.Bytes())
		goul.Log(p.GetLogger(), p.ID, "zlib compress size: %v/%v=%.2f", sizeComp, sizeOrig, float64(sizeComp)/float64(sizeOrig)*100.0)

		out <- &goul.ItemGeneric{Meta: "application/zlib", DATA: b.Bytes()}

		totOrig += int64(sizeOrig)
		totComp += int64(sizeComp)
		count++
	}
	p.SetError(errors.New(goul.ErrPipeInputClosed))
	goul.Log(p.GetLogger(), p.ID, "total %v packets, %v bytes, %.1f%%", count, totOrig, float64(totComp)/float64(totOrig)*100.0)
}

// reverter decompresses items from `in` channel and put it into `out` channel.
func (p *CompressZLib) reverter(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(p.GetLogger(), p.ID, "exit")

	var count, totOrig, totComp int64
	var b bytes.Buffer
	goul.Log(p.GetLogger(), p.ID, "zlib decompressor in looping...")
	for item := range in {
		if item.String() != "application/zlib" {
			goul.Log(p.GetLogger(), p.ID, "item is not a zlib compressed file: %v", item)
			out <- item
			continue
		}

		b.Truncate(0)
		b.Write(item.Data())

		r, err := zlib.NewReader(&b)
		if err != nil {
			goul.Log(p.GetLogger(), p.ID, "could not create zlib reader: %v", err)
			p.SetError(errors.New(CouldNotCreateNewZLibReader))
			continue
		}
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			goul.Log(p.GetLogger(), p.ID, "ioutil error: %v", err)
		}
		r.Close()

		sizeOrig := len(item.Data())
		sizeComp := len(buf)
		goul.Log(p.GetLogger(), p.ID, "zlib dec size: %v/%v", sizeOrig, sizeComp)

		// TODO need to check the type of the buf but... do I deprecate it?
		out <- gopacket.NewPacket(buf, layers.LayerTypeEthernet, gopacket.Default)

		totOrig += int64(sizeOrig)
		totComp += int64(sizeComp)
		count++
	}
	p.SetError(errors.New(goul.ErrPipeInputClosed))
	goul.Log(p.GetLogger(), p.ID, "total %v packets, %v bytes, %.1f%%", count, totOrig, float64(totComp)/float64(totOrig)*100.0)
}
