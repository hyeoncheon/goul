#!/bin/sh
#
# test multiple package by single script for codeclimate.com.
#   https://github.com/codeclimate/test-reporter/blob/master/examples/go_examples.md#example-2

for pkg in `go list ./... | grep -v main`; do
	go test -coverprofile=`echo $pkg | tr / -`.cover $pkg
done
echo "mode: set" > c.out
grep -h -v "^mode:" ./*.cover |sort >> c.out
rm -f *.cover
