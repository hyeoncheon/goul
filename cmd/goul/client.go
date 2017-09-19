package main

import (
	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
)

func client(opts *Options) error {
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	logger := logger(opts)
	router.SetLogger(logger)

	logger.Debugf("initialize device dump on %v...", opts.device)
	reader, err := adapters.NewDevice(opts.device)
	if err != nil {
		logger.Error("couldn't create new device reader: ", err)
		return err
	}
	defer reader.Close()

	if opts.filter != "" {
		logger.Infof("user defined filter: <%v>", opts.filter)
		reader.SetFilter(opts.filter)
	}
	reader.SetOptions(false, 1600, 1)

	logger.Debugf("initialize network connection %v:%v...", opts.addr, opts.port)
	writer, _ := adapters.NewNetwork(opts.addr, opts.port)
	defer writer.Close()

	router.SetReader(reader)
	router.SetWriter(writer)
	//router.AddPipe(&pipes.DebugPipe{ID: "--CP--", Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
	//router.AddPipe(&pipes.DebugPipe{ID: "--CC--", Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})

	control, done, err := router.Run()
	if err != nil {
		logger.Error("couldn't start the router: ", err)
		return err
	}

	if done != nil {
		<-done
		close(control)
	}
	return nil
}
