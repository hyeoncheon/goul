package goul

import (
	"bufio"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// constants
const (
	NetName       = "network adapter"
	ErrCannotDial = "cannot connect to server"
)

// Net is a reader/writer pipe that transport the data via network
type Net struct {
	address string
	err     error
	inLoop  bool
	logger  Logger

	listener net.Listener
	conn     net.Conn
	wBuffer  *bufio.Writer
	rBuffer  *bufio.Reader
}

// NewNetwork creates a new instance of Net and setup network connection.
// If parameter `addr` is not given, it act as a server and make a listener
// as a separated goroutine. Otherwise it connect to the server on given
// address before return.
// Note that, server does not mean a receiver. the client/server role and
// sender/receiver role is orthogonal.
func NewNetwork(addr string, port int) (*Net, error) {
	n := &Net{
		address: addr + ":" + strconv.Itoa(port),
	}

	if addr == "" { // if addr is blank, run as server
		n.logf("in server mode. starting listener on %v...", n.address)
		n.listener, n.err = net.Listen("tcp", n.address)
		if n.listener == nil || n.err != nil {
			return nil, n.err
		}

		go func() {
			for {
				if n.conn != nil {
					time.Sleep(10 * time.Millisecond)
					continue
				}
				if n.listener != nil {
					n.conn, n.err = n.listener.Accept()
					if n.conn == nil {
						n.log("cannot accept: ", n.err)
					} else {
						n.log("client connected from ", n.conn.RemoteAddr())
					}
				}
			}
		}()
	} else {
		n.log("preparing client connection...")
		n.conn, n.err = net.Dial("tcp", n.address)
		if n.conn == nil || n.err != nil {
			return nil, n.err
		}
	}
	return n, nil
}

// Close clean up the resources
func (n *Net) Close() {
	n.log("about to clean up...")
	if n.listener != nil {
		n.log("closing listener...")
		n.listener.Close()
		n.listener = nil
	}
	if n.conn != nil {
		n.log("closing connection...")
		n.conn.Close()
		n.conn = nil
	}
	n.rBuffer = nil
	n.wBuffer = nil
}

// InLoop implements goul.PacketPipe interface
func (n *Net) InLoop() bool {
	return n.inLoop
}

// GetError implements Reader/Writer interface
func (n *Net) GetError() error {
	return n.err
}

// Reader implements Reader interface
func (n *Net) Reader(cmd chan int, out chan Item) {
	defer close(out)
	n.log("reader ready...")

	for {
		select {
		case command := <-cmd:
			switch command {
			case ComInterrupt:
				n.log("shutdown receiving...")
				n.inLoop = false
				n.err = errors.New(ErrPipeInterrupted)
				return
			}
		default: // for unblock from the command channel
		}

		if n.rBuffer != nil {
			// just simple header handling without seek. test purpose.
			header := []byte{}
			for i := 0; i < 2; {
				bt, err := n.rBuffer.ReadByte()
				if err == nil {
					header = append(header, bt)
					i++
				} else {
					n.log("oops! read error ", err)
					n.inLoop = false
					n.err = errors.New(ErrNetworkReadHeader)
					return
				}
			}
			size := binary.BigEndian.Uint16(header)

			data := []byte{}
			remind := int(size)
			for {
				chunk := make([]byte, remind)
				cnt, err := n.rBuffer.Read(chunk)
				if err != nil {
					n.log("ERROR get chunk: ", err)
					n.err = errors.New(ErrNetworkReadChunk)
					n.inLoop = false
					return
				}
				data = append(data, chunk[0:cnt]...)
				remind -= cnt
				if remind == 0 {
					break
				}
			}
			n.logf("%v read %v", remind, len(data))

			// convert raw data into pcap packet if possible.
			// but how and what can I do for compressed data?
			// am I need autodetection? or just let pipeline do it?
			p := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
			if packet, ok := p.(gopacket.Packet); ok {
				out <- packet
			} else {
				out <- &ItemGeneric{Meta: ItemTypeUnknown, DATA: data}
			}
		} else if n.conn != nil {
			n.rBuffer = bufio.NewReader(n.conn)
			n.inLoop = true
			n.log("reading started!")
		} else {
			time.Sleep(50 * time.Millisecond) // no conn but wait. WHY?
		}
	}
}

// Writer implements Writer interface
func (n *Net) Writer(in chan Item) {
	n.log("writer ready...")

	header := make([]byte, 2)
	for {
		select {
		case item, ok := <-in:
			if !ok {
				n.inLoop = false
				n.err = errors.New(ErrPipeInputClosed)
				n.log("could not read channel. writer exit")
				return
			}
			if n.wBuffer != nil {
				data := item.Data()
				binary.BigEndian.PutUint16(header, uint16(len(data)))

				if _, err := n.wBuffer.Write(header); err != nil {
					n.resetConnection(err)
				} else if _, err := n.wBuffer.Write(data); err != nil {
					n.resetConnection(err)
				} else if err := n.wBuffer.Flush(); err != nil {
					n.resetConnection(err)
				} else {
					n.log("sent ", len(data))
				}
			} else if n.conn != nil {
				n.wBuffer = bufio.NewWriter(n.conn)
				n.inLoop = true
				n.log("writing started!")
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (n *Net) resetConnection(err error) {
	if err != nil {
		n.log("error: ", err)
	}
	n.log("clean up connection status...")
	n.conn.Close()
	n.conn = nil
	n.wBuffer = nil
	n.inLoop = false
	n.err = errors.New(ErrNetworkConnectionReset)
	n.log("connection released!")

}

//** logging... -----------------------------------------------------

// SetLogger sets logger for the goul instance.
func (n *Net) SetLogger(l Logger) error {
	n.logger = l
	return nil
}

func (n *Net) log(args ...interface{}) {
	if n.logger != nil {
		args = append([]interface{}{NetName + ": "}, args...)
		n.logger.Debug(args...)
	}
}

func (n *Net) logf(format string, args ...interface{}) {
	if n.logger != nil {
		n.logger.Debugf(NetName+": "+format, args...)
	}
}
