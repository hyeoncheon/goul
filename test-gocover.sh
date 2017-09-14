#!/bin/sh
# vim: set ts=2 sw=2:
#
#  repository: https://github.com/AlekSi/gocoverutil

set -x

app=github.com/hyeoncheon/goul

gocoverutil -coverprofile=cover.out test -covermode=count $app/... && \
	go tool cover -html=cover.out -o cover.html

