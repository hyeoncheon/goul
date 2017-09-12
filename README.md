# Goul, Virtual Port Mirror over Internet (for Cloud)

[![Go Report Card](https://goreportcard.com/badge/github.com/hyeoncheon/goul)](https://goreportcard.com/report/github.com/hyeoncheon/goul)

Goul(거울; Mirror in English) is a tool for virtual network port mirroring
over the Internet, especially for cloud computing environment.

On legacy infrastructure, with many physical switches, we can use a port
mirror(SPAN) on the switch for monitoring and analysing of network, and
connecting a security appliances. But in cloud computing environment, it
is not easy as legacy and in some cases, it is completely impossible.

This tool is for someone like me who want to mirror some port of the
virtual instances or virtual network appliances.

It is now on development using
Go language
with
gopacket
and
pcap
library.

## Features

* Mirror one network device(port) from virtual instance to remote system.
* Selection of Rx, Tx or Both direction of traffic. (Plan)
* Packet filtering based on pcap library's rule.
* Pipelining for filtering, buffering, compression, deduplication, and more.
* Use TCP/IP for transmission over the Internet.
* Support adaptive mode to reduce the impact of production traffic. (Plan)


## Work Flow

Goul captures packets on the source device with reader function and push it
into processing pipeline.  Processing pipeline is a chain of pipe functions
executed as a goroutine.  Another end of the pipeline is connected to writer
function and it send the data to the receiver side, over the Internet.

The pipe function is a function with input and output channel. The input
data and output data is usally a packet data but some pipe can modify the
packet into different kind of binary. (for example, gzip compressed data)
You can write your own pipe function with this interface for your own
reason. for example, counting specific packets, discard some protocol,
or deduplication of the same packets.

In the receiver side, which is the Goul is running as receiver, the reader
function receives all data from network socket, which is connected to the
sender via Internet. It pushs the data into the reverse pipeline.
The reverse pipeline is also a chain of pipe functions as same as processing
pipeline but reverse order of it. (because the datatype must be matched by
the stacked pipes)

The reverse pipeline recovers the packet into its original form and pass
it to the writer function. Writer function for the receiver mode is packet
injector. After that, the packets on the network interface is virtually
same as original source device. (Simply, "Copied")

![Goul Generic](docs/goul.generic.png)

In above configuration, the packets captured on network `A` will be passed
through the process pipeline on left side, then sended to the right side
over the Internet. In the right side, the receiver reverse it(for example
decompressing it,...) and inject into a local network interface.

The target network interface is connected to the switch port which is
configured as source port of the mirror configuration set. The packet on
the source port will be mirrored by switch configuration and then it passed
into Network Analyzer.

If the network to be monitored is simply, or the analyzer is dedicated to
this set, then you can configure it as below. the only difference is, there
is no switch for mirroring. The Network Analyzer is directly connected to
the receiver server.

![Goul Direct](docs/goul.direct.png)

### Pipeline Details

You can write your own pipe function with input and output channel.
There is no limitation on the job of the function. (Currently, the packet
already contains these functions: zlib and gzip compression and its reverse,
counting, and printing the packet details for debugging.)

The pipe function must be wrote as a pair: compress then decompress as reverse,
encoding then decoding as reverse, and so one. Exception for this rule is a
transparent function which is the output data is same as its input data.
For example, Counter or packet printer is a transparent function.

If the function used in process pipeline is not transparent one, then the
reverse function must be exist on the receiver side.

## Install

Installation of Goul is same as any Go programs. Just get it.
But while compiling it, it needs `libpcap` development packet so you need
to install it before get. Below is the installation process for Ubuntu
Linux. (Or other Debian based Linux distributions)

```console
$ sudo apt-get install libpcap-dev
<...>
$ go get github.com/hyeoncheon/goul/cmd/goul
$ 
$ ls $GOPATH/pkg/linux_amd64/github.com/hyeoncheon
goul  goul.a
$ 
$ ls $GOPATH/bin
goul
$ 
```

## Running

It runs in 2x2 mode.
First, as described above, it can run as sender or receiver. It also run
as server or client. This two type of modes are orthogonal so you can
run it as sender client and receiver server pair or sender server and
receiver client mode.

The reason why I made these somewhat strange and/or confusing pair of
mode is, I consider some network firewall environment. For example,
if your receiver must be located in the network behind the wall, but
your sender is located outside of the wall, you don't need to configure
or ask firewall open to administrator. Just set the receiver as
client. (Currently I am considering removing of this 2x2 mode and just
make receiver as server and sender as client.
Anyway, currently it support this 2x2 modes.)

Command line options are shown below:

```console
$ goul -h
Goul 0.1

Goul is a packet capture program for cloud environment.

If it runs as capturer mode, it captures all packets on local network
interface and sends them to remote receiver over internet.
The other side, while it runs as receiver mode, it receives packets from
remote capturer and inject them into the interface on the system.

Usage: goul [-Dhlrtv] [-a value] [-d value] [-p value] filters ...
 -a, --addr=value  address to connect (for client)
 -D, --debug       debugging mode (print log messages)
 -d, --dev=value   network interface to read/write
 -h, --help        help
 -l, --list        list network devices
 -p, --port=value  address to connect (default is 6001)
 -r, --recv        run as receiver
 -t, --test        test mode (no injection)
 -v, --version     show version of goul
$
```

If you just run this like below, it runs as sender server.

```console
$ sudo ./goul
```

For sender client, you can run Goul like below:

```console
$ sudo ./goul --addr 10.0.0.1
```

10.0.0.1 is a IP address of the receiver server. For serve this sender
client, means receiver server, you can simply run Goul as below:


```console
$ sudo ./goul --recv
```

For receiver client, as you can imagin,

```console
$ sudo ./goul --recv --address 10.0.0.1
```

will be work. for server and client, it uses TCP port 6001 but you can
set your own port number with `-p` flag followed by port number.

The default network device is `eth0` but you can configure it with `-d` flag.
If you don't know which network interface name is for you, you can simply
try `-l` flag for listing all possible interfaces.

```console
$ goul -l

Devices:
* eth0
  - IP address:  10.0.0.2
  - IP address:  fe80::465:0000:0000:40000
* eth1
  - IP address:  192.168.0.2
  - IP address:  fe80::4a8:0000:0000:0000
* any
* lo
  - IP address:  127.0.0.1
  - IP address:  ::1
$ 
```

The Goul is currently on development and it can be unstable. So if you
just want to test it without real packet injection, use `-t` flag.
It allows receiver runs in testing mode then it does not inject but
just display the number of packets it received.

Along with `-t` flag, `-D` flag is useful for testing. It turns Goul in
verbose mode and it print out a lots of messages while running.

Please note that, The Goul's default capturing filter is `ip`. so it
just ignore all non-IP protocols including L2 level protocols.
If you want to set more specific filter, you can put it in the end of
the command line. (`filters ...`) The filter rule is same as other
`pcap` based application like `tcpdump`. So you can set
`port 80 and port 443` as filter for getting `HTTP` and `HTTPS` traffic.


Have fun with packets! and the Goul!


## Current Status

* Support simple capturing and injection of packet over the Internet.
* Currently all compression features are disabled by default.
  * I found that it consumes CPU but the compression ratio is not effective.
* Currently Pipeline configuration from command line is not supported.


## Caution!!!

THIS PROGRAM IS WORKING WITH LOW LEVEL NETWORKING FEATURES. DO NOT USE THIS
PROGRAM IF YOU DO NOT UNDERSTAND WHAT IT DOES.

ESPECIALLY, DO NOT FORWARD L2 MANAGEMENT PROTOCOL LIKE `STP` INTO THE NORMALLY
CONFIGURED PORT. IT CAN BREAK YOUR ENTIRE NETWORK.


## Author

Yonghwan SO https://github.com/sio4, http://www.sauru.so


## Copyright (GNU General Public License v3.0)

Copyright 2016 Yonghwan SO

This program is free software; you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation; either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program; if not, write to the Free Software Foundation, Inc., 51
Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA

