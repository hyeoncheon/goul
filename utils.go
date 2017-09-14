package goul

import (
	"errors"
	"fmt"

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

	fmt.Println(`\nYour system has following device(s).
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
