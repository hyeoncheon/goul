package pipes

import (
	"time"

	"github.com/hyeoncheon/goul"
)

// DummyPipe ...
type DummyPipe struct {
	goul.Pipe
	ID string
}

// Convert implements interface Pipe/Converter
func (p *DummyPipe) Convert(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(p.converter, in, message)
}

// Revert implements interface Pipe/Reverter
func (p *DummyPipe) Revert(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(p.reverter, in, message)
}

// complex, non-blocking loop over the input channel.
func (p *DummyPipe) converter(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(p.GetLogger(), p.ID, "exit")
	goul.Log(p.GetLogger(), p.ID, "converter in looping...")

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
			out <- item
			goul.Log(p.GetLogger(), p.ID, "works %d/%d times", w, i)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// simple, blocking loop over the input channel.
func (p *DummyPipe) reverter(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(p.GetLogger(), p.ID, "exit")
	goul.Log(p.GetLogger(), p.ID, "reverter in looping...")

	i := 0
	for item := range in {
		i++
		out <- item
		goul.Log(p.GetLogger(), p.ID, "works %d times", i)
	}
	goul.Log(p.GetLogger(), p.ID, "channel closed")
}
