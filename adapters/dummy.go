package adapters

import (
	"time"

	"github.com/hyeoncheon/goul"
)

// DummyAdapter ...
type DummyAdapter struct {
	goul.Adapter
	ID string
}

// Read implements interface Adapter
func (a *DummyAdapter) Read(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(a.reader, in, message)
}

// Write implements interface Adapter
func (a *DummyAdapter) Write(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(a.writer, in, message)
}

// complex, non-blocking loop over the input channel.
func (a *DummyAdapter) reader(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(a.GetLogger(), a.ID, "exit")
	goul.Log(a.GetLogger(), a.ID, "reader in looping...")

	i := 0
	w := 0
	for {
		i++
		select {
		case _, ok := <-in:
			if !ok {
				goul.Log(a.GetLogger(), a.ID, "channel closed")
				return
			}
		default:
			w++
			goul.Log(a.GetLogger(), a.ID, "works %d/%d times", w, i)
			out <- &goul.ItemGeneric{DATA: []byte{4}}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// simple, blocking loop over the input channel.
func (a *DummyAdapter) writer(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer goul.Log(a.GetLogger(), a.ID, "exit")

	goul.Log(a.GetLogger(), a.ID, "writer in looping...")
	i := 0
	for range in {
		i++
		goul.Log(a.GetLogger(), a.ID, "works %d times", i)
	}
	goul.Log(a.GetLogger(), a.ID, "channel closed")
	out <- &goul.ItemGeneric{Meta: "message", DATA: []byte("channel closed. done")}
}
