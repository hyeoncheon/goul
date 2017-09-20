package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
)

func run(opts *Options) {
	var err error
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	logger := logger(opts)
	router.SetLogger(logger)

	if opts.isServer {
		logger.Debugf("initialize network connection %v:%v...", opts.addr, opts.port)
		reader, _ := adapters.NewNetwork(opts.addr, opts.port)
		defer reader.Close()

		logger.Debugf("initialize device pump on %v...", opts.device)
		writer, err := adapters.NewDevice(opts.device, opts.isTest)
		if err != nil {
			logger.Error("couldn't create new device reader: ", err)
			return
		}
		defer reader.Close()

		writer.SetOptions(true, 1600, 1)

		router.SetReader(reader)
		router.SetWriter(writer)
		//router.AddPipe(&pipes.CompressZLib{Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})
		//router.AddPipe(&pipes.DebugPipe{Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
	} else {
		logger.Debugf("initialize device dump on %v...", opts.device)
		reader, err := adapters.NewDevice(opts.device, opts.isTest)
		if err != nil {
			logger.Error("couldn't create new device reader: ", err)
			os.Exit(2)
		}
		defer reader.Close()

		if opts.filter != "" {
			logger.Infof("user defined filter: <%v>", opts.filter)
			reader.SetFilter(opts.filter)
		}
		reader.SetOptions(true, 1600, 1)

		logger.Debugf("initialize network connection %v:%v...", opts.addr, opts.port)
		writer, _ := adapters.NewNetwork(opts.addr, opts.port)
		defer writer.Close()

		router.SetReader(reader)
		router.SetWriter(writer)
		//router.AddPipe(&pipes.DebugPipe{Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
		//router.AddPipe(&pipes.CompressZLib{Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
		//router.AddPipe(&pipes.DebugPipe{Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})
	}
	control, done, err := router.Run()
	if err != nil {
		logger.Error("couldn't start the router: ", err)
		os.Exit(2)
	}

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
				select {
				case <-control:
				default: // if channel is alive
					close(control)
				}
			}
		}
	}()

	if done != nil {
		<-done
		select {
		case <-control:
		default: // if channel is alive
			close(control)
		}
	}
}

//** utilities...

func logger(opts *Options) goul.Logger {
	if opts.isDebug {
		return goul.NewLogger("debug")
	}
	return goul.NewLogger("info")
}
