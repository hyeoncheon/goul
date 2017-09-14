package pipes_test

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
)

func Test_PacketPrinter(t *testing.T) {
	pt := &TransparentPipeTest{PacketPipe: &pipes.PacketPrinter{}, T: t}
	pt.Run()
	pt.RawToPacket()
}

func Test_DataCounter(t *testing.T) {
	pt := &TransparentPipeTest{PacketPipe: &pipes.DataCounter{}, T: t}
	pt.Run()
}

func Test_NullWriter(t *testing.T) {
	pt := &WriterTest{Writer: &pipes.NullWriter{}, T: t}
	pt.Run()
}

//** writer pipe test -----------------------------------------------

type WriterTest struct {
	Writer goul.Writer
	T      *testing.T
}

func (p *WriterTest) Run() {
	p.Normal()
}

func (p *WriterTest) Normal() {
	r := require.New(p.T)
	packet, err := SetupPacket("test message")
	r.NoError(err)

	p.Writer.SetLogger(goul.NewLogger("debug"))
	in := make(chan goul.Item)
	go p.Writer.Writer(in)

	in <- packet
	in <- packet
	time.Sleep(200 * time.Millisecond)
	r.Nil(p.Writer.GetError())

	close(in)
	for p.Writer.InLoop() {
		time.Sleep(100 * time.Millisecond)
	}
	r.NotNil(p.Writer.GetError())
	r.Equal(goul.ErrPipeInputClosed, p.Writer.GetError().Error())
}

//** packet pipe test -----------------------------------------------

type PacketPipeTest struct {
	PacketPipe goul.PacketPipe
	T          *testing.T
}

func (p *PacketPipeTest) Run() {
	p.Normal()
	p.CompressedToPipe()
	p.PacketToReverse()
}

func (p *PacketPipeTest) Normal() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := SetupPacket(data)
	r.NoError(err)

	err = CheckPacket(packet, data)
	r.NoError(err, "CheckPacket: %v", err)

	in := make(chan goul.Item)
	out := make(chan goul.Item)
	mid := make(chan goul.Item)

	go p.PacketPipe.Pipe(in, mid)
	go p.PacketPipe.Reverse(mid, out)

	in <- packet
	result := <-out
	err = CheckPacket(result, data)
	r.NoError(err, "CheckPacket: %v", err)

	r.Nil(p.PacketPipe.GetError())
	close(in) // channel mid will be closed automatically by Pipe.
	for p.PacketPipe.InLoop() {
		time.Sleep(100 * time.Millisecond)
	}
	r.NotNil(p.PacketPipe.GetError())
	r.Equal(goul.ErrPipeInputClosed, p.PacketPipe.GetError().Error())
}

func (p *PacketPipeTest) CompressedToPipe() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := SetupPacket(data)
	r.NoError(err)

	in1 := make(chan goul.Item)
	out1 := make(chan goul.Item)
	in2 := make(chan goul.Item)
	out2 := make(chan goul.Item)

	go p.PacketPipe.Pipe(in1, out1)
	go p.PacketPipe.Reverse(in2, out2)

	in1 <- packet
	gz1 := <-out1
	err = CheckPacket(gz1, data)
	r.Error(err)
	r.Contains(err.Error(), "NoPacket")
	r.Contains(gz1.String(), "application")

	in1 <- gz1
	gz2 := <-out1
	r.Equal(gz1, gz2, "output of re-enterance must be same")

	in2 <- gz2
	pkt := <-out2
	err = CheckPacket(pkt, data)
	r.NoError(err, "CheckPacket: %v", err)

	close(in1)
	close(in2)
}

func (p *PacketPipeTest) PacketToReverse() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := SetupPacket(data)
	r.NoError(err)

	in := make(chan goul.Item)
	out := make(chan goul.Item)

	go p.PacketPipe.Reverse(in, out)

	in <- packet
	pkt := <-out
	err = CheckPacket(pkt, data)
	r.NoError(err, "CheckPacket: %v", err)

	//// NOTE: DO NOT TEST Error() AGAIN. IT DO NOT RESET EVERY RUN.
	close(in)
}

//** transparent pipe test ------------------------------------------

type TransparentPipeTest struct {
	PacketPipe goul.PacketPipe
	T          *testing.T
}

func (p *TransparentPipeTest) Run() {
	p.NormalDirection()
	p.ReverseDirection()
}

func (p *TransparentPipeTest) NormalDirection() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := SetupPacket(data)
	r.NoError(err)

	err = CheckPacket(packet, data)
	r.NoError(err, "CheckPacket: %v", err)

	in := make(chan goul.Item)
	out := make(chan goul.Item)
	mid := make(chan goul.Item)

	go p.PacketPipe.Pipe(in, mid)
	go p.PacketPipe.Reverse(mid, out)

	in <- packet
	result := <-out
	err = CheckPacket(result, data)
	r.NoError(err, "CheckPacket: %v", err)

	r.Nil(p.PacketPipe.GetError())
	close(in) // channel mid will be closed automatically by Pipe.
	for p.PacketPipe.InLoop() {
		time.Sleep(100 * time.Millisecond)
	}
	r.NotNil(p.PacketPipe.GetError())
	r.Equal(goul.ErrPipeInputClosed, p.PacketPipe.GetError().Error())
}

func (p *TransparentPipeTest) ReverseDirection() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := SetupPacket(data)
	r.NoError(err)

	err = CheckPacket(packet, data)
	r.NoError(err, "CheckPacket: %v", err)

	in := make(chan goul.Item)
	out := make(chan goul.Item)
	mid := make(chan goul.Item)

	go p.PacketPipe.Reverse(in, mid)
	go p.PacketPipe.Pipe(mid, out)

	in <- packet
	result := <-out
	err = CheckPacket(result, data)
	r.NoError(err, "CheckPacket: %v", err)

	//// NOTE: DO NOT TEST Error() AGAIN. IT DO NOT RESET EVERY RUN.
	close(in) // channel mid will be closed automatically by Pipe.
}

func (p *TransparentPipeTest) RawToPacket() {
	r := require.New(p.T)

	data := "Test String"
	packet, err := SetupPacket(data)
	r.NoError(err)

	err = CheckPacket(packet, data)
	r.NoError(err, "CheckPacket: %v", err)

	in := make(chan goul.Item)
	out := make(chan goul.Item)
	mid := make(chan goul.Item)

	go p.PacketPipe.Pipe(in, mid)
	go p.PacketPipe.Reverse(mid, out)

	// packet in, packet out
	in <- packet
	result := <-out
	err = CheckPacket(result, data)
	r.NoError(err, "CheckPacket: %v", err)

	// raw packet in, packet out
	in <- &goul.ItemGeneric{Meta: goul.ItemTypeRawPacket, DATA: packet.Data()}
	result = <-out
	err = CheckPacket(result, data)
	r.NoError(err, "CheckPacket: %v", err)

	// raw packet without header in, raw packet out
	in <- &goul.ItemGeneric{Meta: "", DATA: packet.Data()}
	result = <-out
	err = CheckRawPacket(result, data)
	r.NoError(err, "CheckPacket: %v", err)

	//// NOTE: DO NOT TEST Error() AGAIN. IT DO NOT RESET EVERY RUN.
	close(in) // channel mid will be closed automatically by Pipe.
}

//*** utilities... --------------------------------------------------

func CheckPacket(item goul.Item, data string) error {
	p, ok := item.(gopacket.Packet)
	if !ok {
		return errors.New("NoPacket")
	}
	if L := p.LinkLayer(); L == nil || L.LayerType() != layers.LayerTypeEthernet {
		return errors.New("NotEthernet")
	}
	if L := p.NetworkLayer(); L == nil || L.LayerType() != layers.LayerTypeIPv4 {
		return errors.New("NotIPv4")
	}
	if L := p.TransportLayer(); L == nil || L.LayerType() != layers.LayerTypeTCP {
		return errors.New("NotTCP")
	}
	if L := p.ApplicationLayer(); L == nil || string(L.Payload()) != data {
		return errors.New("DataMismatch")
	}
	return nil
}

func CheckRawPacket(item goul.Item, data string) error {
	if item.String() != goul.ItemTypeRawPacket {
		fmt.Println("meta mismatch: ", item.String())
	}

	packet := gopacket.NewPacket(item.Data(), layers.LayerTypeEthernet, gopacket.Default)
	p, ok := packet.(gopacket.Packet)
	if !ok {
		return errors.New("NoPacket")
	}
	if L := p.LinkLayer(); L == nil || L.LayerType() != layers.LayerTypeEthernet {
		return errors.New("NotEthernet")
	}
	if L := p.NetworkLayer(); L == nil || L.LayerType() != layers.LayerTypeIPv4 {
		return errors.New("NotIPv4")
	}
	if L := p.TransportLayer(); L == nil || L.LayerType() != layers.LayerTypeTCP {
		return errors.New("NotTCP")
	}
	if L := p.ApplicationLayer(); L == nil || string(L.Payload()) != data {
		return errors.New("DataMismatch")
	}
	return nil
}

func SetupPacket(data string) (gopacket.Packet, error) {
	var err error

	tcp := layers.TCP{
		SrcPort: layers.TCPPort(1234),
		DstPort: layers.TCPPort(80),
		Seq:     11050,
		Ack:     0,
		RST:     true,
	}
	ip := layers.IPv4{
		Protocol: layers.IPProtocolTCP, // must to decode next
		SrcIP:    net.IP{127, 0, 0, 1},
		DstIP:    net.IP{8, 8, 8, 8},
		Version:  4,
	}
	ether := layers.Ethernet{
		EthernetType: layers.EthernetTypeIPv4, // must to decode next
		SrcMAC:       net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC:       net.HardwareAddr{0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD},
	}
	if err := tcp.SetNetworkLayerForChecksum(&ip); err != nil { // must
		return nil, errors.New("SetNetworkLayerForChecksum")
	}

	buf := gopacket.NewSerializeBuffer()
	opt := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	err = gopacket.SerializeLayers(buf, opt,
		&ether,
		&ip,
		&tcp,
		gopacket.Payload([]byte(data)),
	)
	if err != nil {
		return nil, errors.New("SerializeLayers")
	}

	rawPacket := buf.Bytes()
	return gopacket.NewPacket(rawPacket, layers.LayerTypeEthernet, gopacket.Default), nil
}
