package pipes_test

import (
	"testing"

	"github.com/google/gopacket"
	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
	. "github.com/hyeoncheon/goul/tester_test"
)

func Test_Debugger_1_All(t *testing.T) {
	r := require.New(t)
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	router.SetLogger(goul.NewLogger("debug"))
	router.SetReader(&GeneratorAdapter{ID: "G-----", Adapter: &goul.BaseAdapter{}})
	router.SetWriter(&GeneratorAdapter{ID: "-----W", Adapter: &goul.BaseAdapter{}})
	router.AddPipe(&pipes.DebugPipe{ID: "-P----", Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
	router.AddPipe(&pipes.DebugPipe{ID: "----C-", Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})

	control, done, err := router.Run()
	r.NoError(err)
	r.NotNil(control)
	r.NotNil(done)

	// generate pcap packet
	control <- &goul.ItemGeneric{Meta: "packet", DATA: []byte("TestData")}
	out := <-done
	p, ok := out.(gopacket.Packet)
	r.True(ok)
	r.NotNil(p)
	r.NotNil(p.ApplicationLayer())
	r.Equal("TestData", string(p.ApplicationLayer().Payload()))

	// generate raw packet
	control <- &goul.ItemGeneric{Meta: "rawpacket", DATA: []byte("TestData")}
	out = <-done
	p, ok = out.(gopacket.Packet)
	r.True(ok)
	r.NotNil(p)
	r.NotNil(p.ApplicationLayer())
	r.Equal("TestData", string(p.ApplicationLayer().Payload()))

	close(control)
	<-done
}
