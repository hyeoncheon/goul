package adapters

import (
	"bufio"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/hyeoncheon/goul"
)

// constants...
const (
	ErrNetworkWriterNotSupported = "writer not supported for server"
	ErrNetworkReaderNotSupported = "reader not supported for client"
	ErrNetworkReadHeader         = "could not read header from network"
	ErrNetworkReadChunk          = "could not read chunk from network"
)

// NetworkAdapter is normal mode networking adapter.
type NetworkAdapter struct {
	goul.Adapter
	ID       string
	err      error
	address  string
	isServer bool
	listener *net.TCPListener
}

// Read implements interface Adapter
func (a *NetworkAdapter) Read(ctrl chan goul.Item, message goul.Message) (chan goul.Item, error) {
	out := make(chan goul.Item, goul.ChannelSize)
	if a.isServer {
		laddr, _ := net.ResolveTCPAddr("tcp", a.address)
		a.listener, a.err = net.ListenTCP("tcp", laddr)
		if a.err != nil {
			return nil, a.err
		}
		go a.listen(ctrl, out)
	} else {
		return nil, errors.New(ErrNetworkReaderNotSupported)
	}
	return out, nil
}

// Write implements interface Adapter
func (a *NetworkAdapter) Write(in chan goul.Item, message goul.Message) (chan goul.Item, error) {
	done := make(chan goul.Item)
	if !a.isServer {
		conn, err := a.connect()
		if err != nil {
			a.err = err
			return nil, a.err
		}
		go a.writer(in, done, conn)
	} else {
		return nil, errors.New(ErrNetworkWriterNotSupported)
	}
	return done, nil
}

// complex, non-blocking loop over the input channel.
func (a *NetworkAdapter) reader(ctrl, out chan goul.Item, conn *net.TCPConn) {
	defer goul.Log(a.GetLogger(), a.ID+"-rcv", "exit")
	defer conn.Close()

	goul.Log(a.GetLogger(), a.ID+"-rcv", "reader in looping...")
	buffer := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	var i int
	var header [2]byte
	var err error
	var nerr net.Error
	var ok bool
	var size uint16
	var remind int
	var data []byte
	var cnt int
	var packet gopacket.Packet
	for {
		for i = 0; i < 2; {
			header[i], err = buffer.ReadByte()
			if err == nil {
				i++
			} else {
				if nerr, ok = err.(net.Error); ok && nerr.Timeout() {
					conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
					continue
				} else {
					a.SetError(errors.New(ErrNetworkReadHeader))
					goul.Log(a.GetLogger(), a.ID+"-rcv", "oops! couldn't read header: %v", err)
					return
				}
			}
		}
		size = binary.BigEndian.Uint16(header[:])

		data = []byte{}
		remind = int(size)
		for {
			chunk := make([]byte, remind)
			cnt, err = buffer.Read(chunk)
			if err != nil {
				if nerr, ok = err.(net.Error); ok && nerr.Timeout() {
					conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
					continue
				} else {
					a.SetError(errors.New(ErrNetworkReadChunk))
					goul.Log(a.GetLogger(), a.ID+"-rcv", "oops! couldn't read chunk: %v", err)
					return
				}
			}
			data = append(data, chunk[0:cnt]...)
			remind -= cnt
			if remind == 0 {
				break
			}
		}
		goul.Log(a.GetLogger(), a.ID+"-rcv", "read %v %v", len(data), remind)

		select {
		case _, ok = <-ctrl:
			if !ok {
				goul.Log(a.GetLogger(), a.ID+"-rcv", "channel closed b4 write")
				return
			}
		default:
		}
		// convert raw data into pcap packet if possible.
		// but how and what can I do for compressed data?
		// am I need autodetection? or just let pipeline do it?
		// TODO: This code does not work properly. gzip also treated as packet.
		// TODO: Please add mime type in header or just remove all pipes.
		packet = gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
		if packet != nil {
			out <- packet
		} else {
			out <- &goul.ItemGeneric{Meta: goul.ItemTypeUnknown, DATA: data}
		}
	}
}

// writer is function for client module.
func (a *NetworkAdapter) writer(in, done chan goul.Item, conn net.Conn) {
	defer close(done) //! if it runs on server?
	defer goul.Log(a.GetLogger(), a.ID+"-snd", "exit")
	defer conn.Close()

	// preparing write buffers
	header := make([]byte, 2)
	buffer := bufio.NewWriter(conn)

	goul.Log(a.GetLogger(), a.ID+"-snd", "writer in looping...")
	for item := range in {
		data := item.Data()
		size := len(data)
		binary.BigEndian.PutUint16(header, uint16(size))

		hc, err1 := buffer.Write(header)
		bc, err2 := buffer.Write(data)
		err3 := buffer.Flush()
		if err1 != nil || err2 != nil || err3 != nil {
			goul.Log(a.GetLogger(), a.ID+"-snd", "oops! couldn't write to: %v/%v/%v", err1, err2, err3)
			a.SetError(errors.New("ErrNetAdapterWriteError"))
			//! return or signal to the parent?
			return
		}
		if hc != 2 || bc != size {
			goul.Log(a.GetLogger(), a.ID+"-snd", "oops! couldn't fully write: H:%v / B:%v (E:%v)", hc, bc, size)
			a.SetError(errors.New("ErrNetAdapterSizeNotMatch"))
			return
		}
		goul.Log(a.GetLogger(), a.ID+"-snd", "sent %v", size)
	}
	goul.Log(a.GetLogger(), a.ID+"-snd", "channel closed")
	done <- goul.Messages["closed"]
}

// NewNetwork ...
func NewNetwork(addr string, port int) (*NetworkAdapter, error) {
	a := &NetworkAdapter{
		Adapter:  &goul.BaseAdapter{},
		ID:       "net",
		address:  addr + ":" + strconv.Itoa(port),
		isServer: addr == "",
	}
	return a, nil
}

// Close implements Adapter:
func (a *NetworkAdapter) Close() error {
	goul.Log(a.GetLogger(), a.ID, "cleanup...")
	if a.listener != nil {
		a.listener.Close()
	}
	return nil
}

func (a *NetworkAdapter) connect() (net.Conn, error) {
	goul.Log(a.GetLogger(), a.ID, "preparing client connection...")
	conn, err := net.Dial("tcp", a.address)
	return conn, err
}

func (a *NetworkAdapter) listen(in, out chan goul.Item) {
	defer close(out)
	defer goul.Log(a.GetLogger(), a.ID+"-listener", "exit")
	defer a.listener.Close()

	goul.Log(a.GetLogger(), a.ID+"-listener", "preparing listener...")
	for {
		select {
		case _, ok := <-in:
			if !ok {
				goul.Log(a.GetLogger(), a.ID+"-listener", "channel closed")
				return
			}
		default:
			a.listener.SetDeadline(time.Now().Add(1 * time.Second))
			conn, err := a.listener.AcceptTCP()
			if err == nil {
				goul.Log(a.GetLogger(), a.ID+"-listener", "connected from %v", conn.RemoteAddr())
				go a.reader(in, out, conn)
			} else {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					continue
				} else {
					goul.Log(a.GetLogger(), a.ID+"-listener", "couldn't accept: %v", err)
					return
				}
			}
		}
	}
}
