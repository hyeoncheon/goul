package main

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_RunServer(t *testing.T) {
	r := require.New(t)

	svrOpts := &Options{
		isDebug:  true,
		isTest:   false,
		isServer: true,
		port:     6060,
		device:   "bond9", // does not exist
		filter:   "port 80",
	}

	err := run(svrOpts)
	r.Error(err)
	//r.EqualError(err, ErrCouldNotStartTheRouter) // in my local (libpcap 1.9.1)
	//r.EqualError(err, ErrCouldNotCreateDeviceReader) // in travis-ci (libpcap 1.7.4)
}

func Test_RunTestServer(t *testing.T) {
	r := require.New(t)

	svrOpts := &Options{
		isDebug:  true,
		isTest:   true,
		isServer: true,
		port:     6099,
		device:   "bond9",
		filter:   "port 80",
	}

	//*** testing for singla handling...
	wg := sync.WaitGroup{}
	sig := make(chan os.Signal, 1)
	var goerr error
	wg.Add(1)
	go func() {
		goerr = run(svrOpts, sig)
		wg.Done()
	}()
	time.Sleep(1 * time.Second)
	sig <- syscall.SIGUSR1
	sig <- syscall.SIGINT
	wg.Wait()
	r.NoError(goerr)
}

func Test_RunClient(t *testing.T) {
	r := require.New(t)

	svrOpts := &Options{
		isDebug:  true,
		isTest:   false,
		isServer: false,
		addr:     "localhost",
		port:     6060,
		device:   "lo",
		filter:   "port 80",
	}

	err := run(svrOpts)
	r.EqualError(err, ErrCouldNotStartTheRouter) // permission
}

func Test_Logger(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		isDebug: true,
	}

	loggerDebug := logger(opts)
	r.NotNil(loggerDebug)

	loggerNormal := logger(&Options{})
	r.NotNil(loggerNormal)
}
