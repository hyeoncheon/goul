package adapters

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"github.com/hyeoncheon/goul"
)

// constants...
const (
	defaultID          = "PcapRD"
	defaultSnapLen     = 1600
	defaultPromiscuous = false
	defaultTimeout     = 1
	defaultFilter      = "ip"

	ErrDeviceAdapterNotInitialized = "DeviceAdapterNotInitialized"
	ErrCouldNotActivate            = "CouldNotActivate"
)

// DeviceAdapter is an adapter for the network device interfacing.
// This is the most important adapter of Goul. It is used as a reader
// adapter for the sender and a writer adapter for the receiver.
//
// Please note that the adapter MUST be initialized with NewDevice()
// function so that it can be initialized with initialization of device
// and automatically inherit the BaseAdapter that implements underlying
// CommonMixin. Otherwise, it does not work properly.
type DeviceAdapter struct {
	goul.Adapter
	ID          string
	err         error
	device      string
	snaplen     int
	promiscuous bool
	timeout     time.Duration
	filter      string

	handle         *pcap.Handle
	inactiveHandle *pcap.InactiveHandle
}

// Read implements interface Adapter
func (a *DeviceAdapter) Read(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(a.reader, in, message)
}

// Write implements interface Adapter
func (a *DeviceAdapter) Write(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	return goul.Launch(a.writer, in, message)
}

// reader read packets from device and push it into output channel.
func (a *DeviceAdapter) reader(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer a.recover()
	defer goul.Log(a.GetLogger(), a.ID, "exit")

	a.err = a.activate()
	if a.err != nil {
		a.SetError(a.err)
		goul.Error(a.GetLogger(), a.ID, "activation error: %v", a.err)
		return
	}

	if a.err = a.handle.SetBPFFilter(a.filter); a.err != nil {
		a.err = errors.New(ErrCouldNotActivate + " " + a.err.Error())
		return
	}
	packetSource := gopacket.NewPacketSource(a.handle, a.handle.LinkType())
	packetChannel := packetSource.Packets()

	goul.Log(a.GetLogger(), a.ID, "capturing in looping...")
	for {
		select {
		case _, ok := <-in:
			if !ok {
				goul.Log(a.GetLogger(), a.ID, "channel closed")
				return
			}
		case packet := <-packetChannel:
			out <- packet
		default: // for non-blocking looping
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// writer write out the packets from input channel
func (a *DeviceAdapter) writer(in, out chan goul.Item, message goul.Message) {
	defer close(out)
	defer a.recover()
	defer goul.Log(a.GetLogger(), a.ID, "exit")

	a.err = a.activate()
	if a.err != nil {
		a.SetError(a.err)
		goul.Error(a.GetLogger(), a.ID, "activation error: %v", a.err)
		return
	}

	goul.Log(a.GetLogger(), a.ID, "writer in looping...")
	for item := range in {
		if p, ok := item.(gopacket.Packet); ok {
			a.handle.WritePacketData(p.Data())
		}
	}
	goul.Log(a.GetLogger(), a.ID, "channel closed")
	out <- &goul.ItemGeneric{Meta: "message", DATA: []byte("channel closed. done")}
}

// NewDevice returns new device adapter.
func NewDevice(dev string) (*DeviceAdapter, error) {
	a := &DeviceAdapter{
		ID:          defaultID,
		device:      dev,
		snaplen:     defaultSnapLen,
		promiscuous: defaultPromiscuous,
		timeout:     defaultTimeout,
		filter:      defaultFilter,
		Adapter:     &goul.BaseAdapter{},
	}
	a.inactiveHandle, a.err = pcap.NewInactiveHandle(a.device)
	return a, a.err
}

// Close clean up resources on device adapter.
func (a *DeviceAdapter) Close() error {
	defer a.recover()
	goul.Log(a.GetLogger(), a.ID, "cleanup...")
	if a.handle != nil {
		a.handle.Close()
	}
	if a.inactiveHandle != nil {
		a.inactiveHandle.CleanUp()
	}
	return nil
}

// SetOptions sets capture options to inactive handler.
func (a *DeviceAdapter) SetOptions(promisc bool, len int, timeout time.Duration) error {
	defer a.recover()

	if a.inactiveHandle == nil {
		a.err = errors.New(ErrDeviceAdapterNotInitialized)
		return a.err
	}
	if a.err = a.inactiveHandle.SetTimeout(timeout * time.Second); a.err != nil {
		goul.Error(a.GetLogger(), a.ID, "set timeout error: %v", a.err)
	} else if a.err = a.inactiveHandle.SetSnapLen(len); a.err != nil {
		goul.Error(a.GetLogger(), a.ID, "set snaplen error: %v", a.err)
	} else if a.err = a.inactiveHandle.SetPromisc(promisc); a.err != nil {
		goul.Error(a.GetLogger(), a.ID, "set promisc error: %v", a.err)
	}
	a.promiscuous = promisc
	a.snaplen = len
	a.timeout = timeout
	return a.err
}

// SetFilter sets filter string which is applied while capturing.
func (a *DeviceAdapter) SetFilter(filter string) error {
	a.filter = filter
	return nil
}

func (a *DeviceAdapter) activate() error {
	defer a.recover()

	if a.inactiveHandle == nil {
		a.err = errors.New(ErrDeviceAdapterNotInitialized)
		return a.err
	}
	if a.handle == nil && a.inactiveHandle != nil {
		a.handle, a.err = a.inactiveHandle.Activate()
		if a.err != nil {
			a.err = errors.New(ErrCouldNotActivate + " " + a.err.Error())
			return a.err
		}
	}
	goul.Log(a.GetLogger(), a.ID, "handle initiated: %v", a.handle)
	return nil
}

func (a *DeviceAdapter) recover() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "DeviceAdapter recovered from panic!\n")
		fmt.Fprintf(os.Stderr, "Check if NewDevice() is used for creation.\n")
		fmt.Fprintf(os.Stderr, "panic: %v\n", r)
	}
}
