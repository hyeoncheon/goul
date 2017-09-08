package goul

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

// PrintDevices print out all interfaces on the system.
func PrintDevices() {
	fmt.Println("\nDevices:")
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) == 0 {
		fmt.Println(" No compatible devices are found!")
	}

	for _, device := range devices {
		fmt.Println("*", device.Name)
		for _, address := range device.Addresses {
			fmt.Println("  - IP address: ", address.IP)
		}
	}
}
