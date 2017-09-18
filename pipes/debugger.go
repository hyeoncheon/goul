package pipes

import (
	"fmt"
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
func (p *DebugPipe) Convert(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(p.converter, in, message)
}

// Revert implements interface Pipe/Reverter
func (p *DebugPipe) Revert(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
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
				goul.Log(p.GetLogger(), p.ID, "channel closed")
				return
			}
			w++

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
	goul.Log(p.GetLogger(), p.ID, "channel closed")
	goul.Log(p.GetLogger(), p.ID, "counts %d items", i)
}
