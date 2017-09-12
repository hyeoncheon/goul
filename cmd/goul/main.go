package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	getopt "github.com/pborman/getopt/v2"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/transport"
)

// Constants
const (
	PROGRAM = "Goul"
	VERSION = "0.1"
	PORT    = 6001
)

// Options is a structure for running configuration
type Options struct {
	isTest     bool
	isDebug    bool
	isReceiver bool
	addr       string
	port       int
	device     string
	filter     string
	logger     goul.Logger
}

func main() {
	//* initiate with command line arguments...
	opts := getOptions()
	if opts == nil {
		os.Exit(1)
	}

	logLevel := "info"
	if opts.isDebug {
		logLevel = "debug"
	}
	logger := goul.NewLogger(logLevel)

	if opts.filter != "" {
		logger.Infof("user defined filter: <%v>", opts.filter)
	}

	chanCmd := make(chan int, 1)
	gl, err := goul.New(opts.device, opts.isReceiver, chanCmd, opts.isDebug)
	if err != nil {
		logger.Error("could not make a goul session! ", err)
	}
	defer gl.Close()

	gl.SetLogger(logger)
	gl.SetOptions(false, 1600, 1)
	if opts.filter != "" {
		gl.SetFilter(opts.filter)
	}

	//* setup network module
	net, err := transport.New(opts.addr, opts.port)
	if net == nil || err != nil {
		logger.Error("could not prepare the network connection! ", err)
		return
	}
	defer net.Close()

	net.SetLogger(logger)

	if opts.isReceiver {
		gl.SetReader(net)
	} else {
		gl.SetWriter(net)
	}

	//* build reader/writer/processor pipeline
	//gl.AddPipe(&pipes.PacketPrinter{})
	//gl.AddPipe(&pipes.CompressGZip{})
	//gl.AddPipe(&pipes.CompressZLib{})
	//gl.AddPipe(&pipes.DataCounter{})
	//gl.SetWriter(&pipes.NullWriter{})

	//* register singnal handlers and command pipiline...
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			s := <-sig
			logger.Debug("signal caught: ", s)
			switch s {
			case syscall.SIGINT:
				logger.Debug("interrupted! exit gracefully...")
				chanCmd <- goul.ComInterrupt
			}
		}
	}()

	if err := gl.Run(); err != nil {
		logger.Error("Error: ", err)
	}
}

//** getopts...

// getOptions return an Options structure storing parse options.
func getOptions() *Options {
	list := false
	help := false
	version := false

	opts := &Options{
		isTest:     false,
		isDebug:    false,
		isReceiver: false,
		addr:       "",
		port:       PORT,
		device:     "eth0",
	}
	getopt.SetParameters("filters ...")
	getopt.FlagLong(&help, "help", 'h', "help")
	getopt.FlagLong(&list, "list", 'l', "list network devices")
	getopt.FlagLong(&opts.isTest, "test", 't', "test mode (no injection)")
	getopt.FlagLong(&opts.isDebug, "debug", 'D', "debugging mode (print log messages)")
	getopt.FlagLong(&opts.isReceiver, "recv", 'r', "run as receiver")
	getopt.FlagLong(&opts.addr, "addr", 'a', "address to connect (for client)")
	getopt.FlagLong(&opts.port, "port", 'p', "address to connect (default is 6001)")
	getopt.FlagLong(&opts.device, "dev", 'd', "network interface to read/write")
	getopt.FlagLong(&version, "version", 'v', "show version of goul")

	getopt.Parse()
	opts.filter = strings.Join(getopt.Args(), " ")

	if version {
		fmt.Println(versionString)
		return nil
	}
	if help {
		fmt.Println(versionString)
		fmt.Println(helpMessage)
		getopt.Usage()
		return nil
	}
	if list {
		goul.PrintDevices()
		return nil
	}
	return opts
}

const versionString = PROGRAM + " " + VERSION

const helpMessage = `
` + PROGRAM + ` is a packet capture program for cloud environment.

If it runs as capturer mode, it captures all packets on local network
interface and sends them to remote receiver over internet.
The other side, while it runs as receiver mode, it receives packets from
remote capturer and inject them into the interface on the system.
`
