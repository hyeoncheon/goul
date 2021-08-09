package pipes_test

import (
	"testing"

	"github.com/google/gopacket"
	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
	. "github.com/hyeoncheon/goul/testing"
)

func Test_Debugger(t *testing.T) {
	pts := &PipeTestSuiteTransparent{
		C: &pipes.DebugPipe{Pipe: &goul.BasePipe{Mode: goul.ModeConverter}},
		R: &pipes.DebugPipe{Pipe: &goul.BasePipe{Mode: goul.ModeReverter}},
		T: t,
	}
	pts.Run()

	ptsda := &PipeTestSuiteDirectAccess{
		C: &pipes.DebugPipe{},
		R: &pipes.DebugPipe{},
		T: t,
	}
	ptsda.Run()
}

//** pipe test suite transparent ------------------------------------

type PipeTestSuite struct {
	C goul.Pipe
	R goul.Pipe
	T *testing.T
}

func (p *PipeTestSuite) Run() {
	p.Flow()
	p.ConvertRevert()
}

func (p *PipeTestSuite) Flow() {
	r := require.New(p.T)

	var router goul.Router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	router.SetLogger(goul.NewLogger("debug"))
	router.SetReader(&GeneratorAdapter{Adapter: &goul.BaseAdapter{}})
	router.SetWriter(&GeneratorAdapter{Adapter: &goul.BaseAdapter{}})
	router.AddPipe(p.C)
	router.AddPipe(p.R)

	control, done, err := router.Run()
	r.NoError(err, "couldn't start router.Run: %v", err)
	r.NotNil(control)
	r.NotNil(done)

	// generate pcap packet
	control <- &goul.ItemGeneric{Meta: "packet", DATA: []byte("TestData")}
	out := <-done
	packet, ok := out.(gopacket.Packet)
	r.True(ok)
	r.NotNil(packet)
	r.NotNil(packet.ApplicationLayer())
	r.Equal("TestData", string(packet.ApplicationLayer().Payload()))

	// generate raw packet
	control <- &goul.ItemGeneric{Meta: "rawpacket", DATA: []byte("TestData")}
	out = <-done
	packet, ok = out.(gopacket.Packet)
	r.True(ok)
	r.NotNil(packet)
	r.NotNil(packet.ApplicationLayer())
	r.Equal("TestData", string(packet.ApplicationLayer().Payload()))

	close(control)
	<-done
}

func (p *PipeTestSuite) ConvertRevert() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := GeneratePacket(data)
	r.NoError(err, "couldn't generate test packet: %v", err)

	in1 := make(chan goul.Item)
	in2 := make(chan goul.Item)

	out1, err := p.C.Convert(in1, nil)
	r.NoError(err, "couldn't start Converter: %v", err)
	out2, err := p.R.Revert(in2, nil)
	r.NoError(err, "couldn't start Converter: %v", err)

	in1 <- packet // normal direction
	samePacket1 := <-out1
	err = CheckPacket(samePacket1, data)
	r.EqualError(err, "NoPacket", "CheckPacket: %v", err)
	r.Contains(samePacket1.String(), "application")

	in1 <- samePacket1 // push processed data again for re-enterance test.
	samePacket2 := <-out1
	r.Equal(samePacket1, samePacket2, "output of re-enterance must be same")

	// check revert function works or not.
	in2 <- samePacket2 // push processed data into revert processor
	pkt := <-out2
	err = CheckPacket(pkt, data)
	r.NoError(err, "CheckPacket: %v", err)

	in2 <- pkt // push reverted data again for re-enterance test.
	samePacket2 = <-out2
	r.Equal(pkt, samePacket2, "output of re-enterance must be same")

	r.Nil(p.C.GetError())
	close(in1) // channel mid will be closed automatically by Pipe.
	<-out1
	r.NotNil(p.C.GetError())
	r.Equal(goul.ErrPipeInputClosed, p.C.GetError().Error())

	close(in2)
}

//** pipe test suite transparent ------------------------------------

type PipeTestSuiteTransparent struct {
	C goul.Pipe
	R goul.Pipe
	T *testing.T
}

func (p *PipeTestSuiteTransparent) Run() {
	p.Flow()
	p.Convert()
	p.Revert()
}

func (p *PipeTestSuiteTransparent) Flow() {
	r := require.New(p.T)

	var router goul.Router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	router.SetLogger(goul.NewLogger("debug"))
	router.SetReader(&GeneratorAdapter{Adapter: &goul.BaseAdapter{}})
	router.SetWriter(&GeneratorAdapter{Adapter: &goul.BaseAdapter{}})
	router.AddPipe(p.C)
	router.AddPipe(p.R)

	control, done, err := router.Run()
	r.NoError(err, "couldn't start router.Run: %v", err)
	r.NotNil(control)
	r.NotNil(done)

	// generate pcap packet
	control <- &goul.ItemGeneric{Meta: "packet", DATA: []byte("TestData")}
	out := <-done
	packet, ok := out.(gopacket.Packet)
	r.True(ok)
	r.NotNil(packet)
	r.NotNil(packet.ApplicationLayer())
	r.Equal("TestData", string(packet.ApplicationLayer().Payload()))

	// generate raw packet
	control <- &goul.ItemGeneric{Meta: "rawpacket", DATA: []byte("TestData")}
	out = <-done
	packet, ok = out.(gopacket.Packet)
	r.True(ok)
	r.NotNil(packet)
	r.NotNil(packet.ApplicationLayer())
	r.Equal("TestData", string(packet.ApplicationLayer().Payload()))

	close(control)
	<-done
}

func (p *PipeTestSuiteTransparent) Convert() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := GeneratePacket(data)
	r.NoError(err, "couldn't generate test packet: %v", err)

	in1 := make(chan goul.Item)
	in2 := make(chan goul.Item) //! remove this before commit

	out1, err := p.C.Convert(in1, nil)
	r.NoError(err, "couldn't start Converter: %v", err)

	in1 <- packet // normal direction
	samePacket1 := <-out1
	err = CheckPacket(samePacket1, data)
	r.NoError(err, "CheckPacket: %v", err)
	/*
		r.Error(err)
		r.Contains(err.Error(), "NoPacket")
		r.Contains(samePacket1.String(), "application")

		// check re-enterance safety of the converted data again.
		in1 <- samePacket1
		samePacket2 := <-out1
		r.Equal(samePacket1, samePacket2, "output of re-enterance must be same")

		// check revert function works or not.
		in2 <- samePacket2
		pkt := <-out2
		err = CheckPacket(pkt, data)
		r.NoError(err, "CheckPacket: %v", err)
	*/
	r.Nil(p.C.GetError())
	close(in1) // channel mid will be closed automatically by Pipe.
	<-out1
	r.NotNil(p.C.GetError())
	r.Equal(goul.ErrPipeInputClosed, p.C.GetError().Error())

	close(in2)
}

func (p *PipeTestSuiteTransparent) Revert() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := GeneratePacket(data)
	r.NoError(err, "couldn't generate test packet: %v", err)

	in1 := make(chan goul.Item)
	out1, err := p.C.Revert(in1, nil)
	r.NoError(err, "couldn't start Reverter: %v", err)

	in1 <- packet // normal direction
	samePacket1 := <-out1
	err = CheckPacket(samePacket1, data)
	r.NoError(err, "CheckPacket: %v", err)

	r.Nil(p.C.GetError())
	close(in1) // channel mid will be closed automatically by Pipe.
	<-out1
	r.NotNil(p.C.GetError())
	r.Equal(goul.ErrPipeInputClosed, p.C.GetError().Error())
}

//** pipe test suite direct access -- for panic checking ------------

type PipeTestSuiteDirectAccess struct {
	C goul.Pipe
	R goul.Pipe
	T *testing.T
}

func (p *PipeTestSuiteDirectAccess) Run() {
	p.PanicCheck()
}

func (p *PipeTestSuiteDirectAccess) PanicCheck() {
	r := require.New(p.T)

	out, err := p.C.Convert(make(chan goul.Item), nil)
	r.Nil(out)
	r.EqualError(err, "panic")

	out, err = p.R.Revert(make(chan goul.Item), nil)
	r.Nil(out)
	r.EqualError(err, "panic")
}
