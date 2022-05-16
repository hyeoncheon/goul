#!/bin/sh
# vim: set ts=2 sw=2:
#
#  repository: https://github.com/AlekSi/gocoverutil

set -x

go test -v -coverprofile=cover.out -covermode=count ./... && \
	go tool cover -html=cover.out -o cover.html

