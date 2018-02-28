package goul

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/google/gopacket/pcap"
)

// PrintDevices print out all interfaces on the system.
func PrintDevices() error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return errors.New("NoDevices")
	}

	fmt.Println(`Your system has following device(s).
Use name of the device with '-d' flag for override default device 'eth0'.
(e.g. '-d bond0')

Devices:`)
	for _, device := range devices {
		fmt.Println("  *", device.Name)
		for _, address := range device.Addresses {
			fmt.Println("    - IP address: ", address.IP)
		}
	}
	return nil
}

// GoID returns goroutine ID.
//	https://blog.sgmansfield.com/2015/12/goroutine-ids/
//	https://groups.google.com/d/msg/golang-nuts/Nt0hVV_nqHE/bwndAYvxAAAJ
//	https://play.golang.org/p/OeEmT_CXyO
func GoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
