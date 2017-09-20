package main

import (
	"fmt"
	"os"
	"strings"

	getopt "github.com/pborman/getopt/v2"

	"github.com/hyeoncheon/goul"
)

// constants...
const (
	PROGRAM = "goul"
	VERSION = "0.2"
	PORT    = 6001
)

// Options is a structure for running configuration
type Options struct {
	isTest   bool
	isDebug  bool
	isServer bool
	addr     string
	port     int
	device   string
	filter   string
}

func main() {
	opts := getOptions()
	if opts == nil {
		os.Exit(0)
	}
	run(opts)
}

//** getopts...

// getOptions return an Options structure storing parse options.
func getOptions() *Options {
	list := false
	help := false
	version := false

	opts := &Options{
		isTest:   false,
		isDebug:  false,
		isServer: false,
		addr:     "",
		port:     PORT,
		device:   "eth0",
	}
	getopt.SetParameters("filters ...")
	getopt.FlagLong(&help, "help", 'h', "help")
	getopt.FlagLong(&list, "list", 'l', "list network devices")
	getopt.FlagLong(&opts.isTest, "test", 'T', "test mode (no injection)")
	getopt.FlagLong(&opts.isDebug, "debug", 'D', "debugging mode (print log messages)")
	getopt.FlagLong(&opts.isServer, "server", 's', "run as receiver")
	getopt.FlagLong(&opts.addr, "addr", 'a', "address to connect (for client)")
	getopt.FlagLong(&opts.port, "port", 'p', "tcp port number (default is 6001)")
	getopt.FlagLong(&opts.device, "dev", 'd', "network interface to read/write")
	getopt.FlagLong(&version, "version", 'v', "show version of goul")

	getopt.Parse()
	opts.filter = strings.Join(getopt.Args(), " ")

	if version {
		fmt.Println(versionString + "-" + buildNumber)
		return nil
	}
	if help {
		fmt.Println(versionString + "-" + buildNumber)
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

var buildNumber = "head"

const versionString = PROGRAM + " " + VERSION

const helpMessage = `
` + PROGRAM + ` is a packet capture program for cloud environment.

If it runs as capturer mode, it captures all packets on local network
interface and sends them to remote receiver over internet.
The other side, while it runs as receiver mode, it receives packets from
remote capturer and inject them into the interface on the system.
`
