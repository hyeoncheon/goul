package goul_test

// used by network_test.go

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/hyeoncheon/goul"
)

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
