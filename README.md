# Goul, Virtual Port Mirror over Internet (for Cloud)

[![Go Report Card](https://goreportcard.com/badge/github.com/hyeoncheon/goul)](https://goreportcard.com/report/github.com/hyeoncheon/goul)

Goul(=Mirror in Korean: 거울) is a tool for virtual network port mirroring
over the Internet, especially for cloud computing environment.

With legacy infrastructure, with many physical switches, we can configure
a mirror port(SPAN) on the switch for network monitoring, analysis, and
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
* Selection of Rx, Tx or Both direction of traffic.
* Packet filtering based on pcap library's rule.
* Pipelining for filtering, buffering, compression, deduplication, and more.
* Use TCP/IP for transmission over the Internet.
* Adaptive mode to reduce the impact of production traffic.


## Current Status

* Just started to develop basic structure.


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

