package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	getopt "github.com/pborman/getopt/v2"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
)

// Constants
const (
	PROGRAM = "Goul"
	VERSION = "0.1"
)

// Options is a structure for running configuration
type Options struct {
	isReceiver bool
	addr       string
	device     string
	args       []string
	logger     goul.Logger
}

func main() {
	//* initiate with command line arguments...
	opts := getOptions()
	if opts == nil {
		fmt.Println("\nannyeong.")
		os.Exit(1)
	}
	logger := goul.NewLogger("info")
	fmt.Printf(" ...and additional arguments: %v\n", opts.args)

	chanCmd := make(chan int, 1)
	gl, err := goul.New(opts.device, opts.isReceiver, chanCmd)
	if err != nil {
		log.Fatal(err)
	}
	defer gl.Close()

	gl.SetLogger(logger)
	gl.SetOptions(false, 1600, 1)

	/*
		gl.AddPipe(pipes.CompressZLib)
		gl.AddPipe(pipes.DecompressZLib)

		gl.AddPipe(pipes.CompressGZip)
		gl.AddPipe(pipes.DecompressGZip)
	*/

	gl.AddPipe(pipes.PacketPrinter)
	gl.AddPipe(pipes.DataCounter)
	gl.SetWriter(pipes.DataWriter)

	//* register singnal handlers and command pipiline...
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		fmt.Println("\nInterrupted! exit gracefully...")
		chanCmd <- goul.SigINT
	}()

	if err := gl.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

//** getopts...

// getOptions return an Options structure storing parse options.
func getOptions() *Options {
	list := false
	help := false

	opts := &Options{
		isReceiver: false,
		addr:       "",
		device:     "eth0",
	}
	getopt.FlagLong(&help, "help", 'h', "help")
	getopt.FlagLong(&list, "list", 'l', "list network devices")
	getopt.FlagLong(&opts.isReceiver, "recv", 'r', "run as receiver")
	getopt.FlagLong(&opts.addr, "conn", 'c', "address to connect (for client)")
	getopt.FlagLong(&opts.device, "dev", 'd', "network interface to read/write")

	getopt.Parse()
	opts.args = getopt.Args()

	if help {
		fmt.Println(Help)
		getopt.Usage()
		return nil
	}
	if list {
		goul.PrintDevices()
		return nil
	}
	return opts
}

// Help is a message for usage screen.
const Help = `
` + PROGRAM + ` is a packet capture program for cloud environment.

If it runs as capturer mode, it captures all packets on local network
interface and sends them to remote receiver over internet.
The other side, while it runs as receiver mode, it receives packets from
remote capturer and inject them into the interface on the system.
`
