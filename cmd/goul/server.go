package main

import (
	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
)

func server(opts *Options) error {
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	logger := logger(opts)
	router.SetLogger(logger)

	logger.Debugf("initialize network connection %v:%v...", opts.addr, opts.port)
	reader, _ := adapters.NewNetwork(opts.addr, opts.port)
	defer reader.Close()

	/*
		logger.Debugf("initialize device dump on %v...", opts.device)
		reader, err := adapters.NewDevice(opts.device)
		if err != nil {
			logger.Error("couldn't create new device reader: ", err)
			return err
		}
		defer reader.Close()
	*/

	router.SetReader(reader)
	router.SetWriter(&adapters.DummyAdapter{ID: "----DW", Adapter: &goul.BaseAdapter{}})
	//router.AddPipe(&pipes.DebugPipe{ID: "--SP--", Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
	//router.AddPipe(&pipes.DebugPipe{ID: "--SC--", Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})

	control, done, err := router.Run()
	if err != nil {
		logger.Error("couldn't start the router:", err)
	}

	if done != nil {
		<-done
		close(control)
	}
	return nil
}
