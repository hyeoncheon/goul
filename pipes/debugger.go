package pipes

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/hyeoncheon/goul"
)

// DebugPipe ...
type DebugPipe struct {
	goul.Pipe
	ID string
}

// Convert implements interface Pipe/Converter
func (p *DebugPipe) Convert(in chan goul.Item, message goul.Message) (out chan goul.Item, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "DebugPipe#Conver recovered from panic!\n")
			fmt.Fprintf(os.Stderr, "Probably an inheritance problem of pipeline instance.\n")
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = errors.New("panic")
		}
	}()

	if p.ID == "" {
		p.ID = "dbg-convert"
	}
	p.SetError(nil)
	return goul.Launch(p.converter, in, message)
}

// Revert implements interface Pipe/Reverter
func (p *DebugPipe) Revert(in chan goul.Item, message goul.Message) (out chan goul.Item, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "DebugPipe#Revert recovered from panic!\n")
			fmt.Fprintf(os.Stderr, "Probably an inheritance problem of pipeline instance.\n")
			fmt.Fprintf(os.Stderr, "panic: %v\n", r)
			err = errors.New("panic")
		}
	}()

	if p.ID == "" {
		p.ID = "dbg-revert"
	}
	p.SetError(nil)
	return goul.Launch(p.reverter, in, message)
}

// complex, non-blocking loop over the input channel.
func (p *DebugPipe) converter(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(p.GetLogger(), p.ID, "exit")
	goul.Log(p.GetLogger(), p.ID, "debugger in looping...")

	i := 0
	w := 0
	for {
		i++
		select {
		case item, ok := <-in:
			if !ok {
				p.SetError(errors.New(goul.ErrPipeInputClosed))
				goul.Log(p.GetLogger(), p.ID, "channel closed")
				return
			}
			w++

			// TODO: need to clean up. looks strange.
			var packet gopacket.Packet
			if packet, ok = item.(gopacket.Packet); !ok &&
				item.String() == goul.ItemTypeRawPacket {
				packet = gopacket.NewPacket(item.Data(), layers.LayerTypeEthernet, gopacket.Default)
				if packet == nil {
					goul.Log(p.GetLogger(), p.ID, "got RawPacket but failed to convert gopacket!")
					continue
				}
			}
			fmt.Println(packet)
			out <- packet

			goul.Log(p.GetLogger(), p.ID, "works %d/%d times", w, i)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// simple, blocking loop over the input channel.
func (p *DebugPipe) reverter(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(p.GetLogger(), p.ID, "exit")
	goul.Log(p.GetLogger(), p.ID, "counter in looping...")

	i := 0
	for item := range in {
		i++
		out <- item
		goul.Log(p.GetLogger(), p.ID, "%4d, size: %v", i, len(item.Data()))
	}
	p.SetError(errors.New(goul.ErrPipeInputClosed))
	goul.Log(p.GetLogger(), p.ID, "channel closed")
	goul.Log(p.GetLogger(), p.ID, "counts %d items", i)
}
