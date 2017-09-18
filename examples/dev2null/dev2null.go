package main

import (
	"fmt"
	"time"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
	"github.com/hyeoncheon/goul/pipes"
)

func main() {
	dev2null()
}

func dev2null() {
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	logger := goul.NewLogger("debug")
	router.SetLogger(logger)

	reader, err := adapters.NewDevice("eth0")
	if err != nil {
		logger.Error("couldn't create new device reader: ", err)
	}
	defer reader.Close()

	router.SetReader(reader)
	router.SetWriter(&adapters.DummyAdapter{ID: "-----W", Adapter: &goul.BaseAdapter{}})
	router.AddPipe(&pipes.DebugPipe{ID: "-P----", Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
	router.AddPipe(&pipes.DebugPipe{ID: "----C-", Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})

	control, done, err := router.Run()
	if err != nil {
		fmt.Println("could not start the router:", err)
	}
	time.Sleep(1000 * time.Millisecond)

	close(control)
	<-done
	return
}
