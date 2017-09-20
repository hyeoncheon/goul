package goul_test

import (
	"testing"
	"time"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
	"github.com/hyeoncheon/goul/pipes"
	"github.com/stretchr/testify/require"
)

func Test_Pipeline_1_Functions(t *testing.T) {
	r := require.New(t)
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	reader := &adapters.DummyAdapter{ID: "R-----", Adapter: &goul.BaseAdapter{}}
	defer reader.Close()
	writer := &adapters.DummyAdapter{ID: "-----W", Adapter: &goul.BaseAdapter{}}
	defer writer.Close()

	router.SetLogger(goul.NewLogger("debug"))
	router.SetReader(reader)
	router.SetWriter(writer)
	router.AddPipe(&pipes.DummyPipe{ID: "-C----", Pipe: &goul.BasePipe{Mode: goul.ModeConverter}})
	router.AddPipe(&pipes.DummyPipe{ID: "----R-", Pipe: &goul.BasePipe{Mode: goul.ModeReverter}})

	control, done, err := router.Run()
	r.NoError(err)
	r.NotNil(control)
	r.NotNil(done)

	time.Sleep(1 * time.Second)
	close(control)
	message := <-done
	r.Equal("message", message.String())
}

func Test_Pipeline_21_ErrorHandlingNoReaderWriter(t *testing.T) {
	r := require.New(t)
	var err error
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	control, done, err := router.Run()
	r.EqualError(err, goul.ErrRouterNoReaderOrWriter)
	r.Nil(control)
	r.Nil(done)
}

func Test_Pipeline_22_ErrorHandlingNoWriter(t *testing.T) {
	r := require.New(t)
	var err error
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	err = router.SetReader(&goul.BaseAdapter{})
	control, done, err := router.Run()
	r.EqualError(err, goul.ErrRouterNoReaderOrWriter)
	r.Nil(control)
	r.Nil(done)
}

func Test_Pipeline_23_ErrorHandlingNoReader(t *testing.T) {
	r := require.New(t)
	var err error
	var router goul.Router
	router = &goul.Pipeline{Router: &goul.BaseRouter{}}

	err = router.SetWriter(&goul.BaseAdapter{})
	control, done, err := router.Run()
	r.EqualError(err, goul.ErrRouterNoReaderOrWriter)
	r.Nil(control)
	r.Nil(done)
}

func Test_Pipeline_24_ErrorHandlingReadNotImpl(t *testing.T) {
	r := require.New(t)
	var err error

	router := &goul.Pipeline{Router: &goul.BaseRouter{}}
	err = router.SetReader(&goul.BaseAdapter{})
	err = router.SetWriter(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	control, done, err := router.Run() // run with unimplemented pipe
	r.EqualError(err, goul.ErrAdapterReadNotImplemented)
	r.Nil(control)
	r.Nil(done)
}

func Test_Pipeline_25_ErrorHandlingWriteNotImpl(t *testing.T) {
	r := require.New(t)
	var err error

	router := &goul.Pipeline{Router: &goul.BaseRouter{}}
	err = router.SetReader(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	err = router.SetWriter(&goul.BaseAdapter{})
	control, done, err := router.Run() // run with unimplemented pipe
	r.EqualError(err, goul.ErrAdapterWriteNotImplemented)
	r.Nil(control)
	r.Nil(done)
}

func Test_Pipeline_26_ErrorHandlingConvertNotImpl(t *testing.T) {
	r := require.New(t)
	var err error

	router := &goul.Pipeline{Router: &goul.BaseRouter{}}
	err = router.SetReader(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	err = router.SetWriter(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	err = router.AddPipe(&goul.BasePipe{Mode: goul.ModeConverter})
	control, done, err := router.Run() // run with unimplemented pipe
	r.EqualError(err, goul.ErrPipeConvertNotImplemented)
	r.Nil(control)
	r.Nil(done)
	r.Equal(1, len(router.GetPipes()))
}

func Test_Pipeline_27_ErrorHandlingRevertNotImpl(t *testing.T) {
	r := require.New(t)
	var err error

	router := &goul.Pipeline{Router: &goul.BaseRouter{}}
	err = router.SetReader(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	err = router.SetWriter(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	err = router.AddPipe(&goul.BasePipe{Mode: goul.ModeReverter})
	control, done, err := router.Run() // run with unimplemented pipe
	r.EqualError(err, goul.ErrPipeRevertNotImplemented)
	r.Nil(control)
	r.Nil(done)
	r.Equal(1, len(router.GetPipes()))
}

func Test_Pipeline_41_EmptyRouterAddPipe(t *testing.T) {
	r := require.New(t)
	var err error

	router := &goul.Pipeline{}
	err = router.AddPipe(&goul.BasePipe{Mode: goul.ModeReverter})
	r.EqualError(err, "panic")
}

func Test_Pipeline_42_EmptyRouterRun(t *testing.T) {
	r := require.New(t)
	var err error

	router := &goul.Pipeline{}
	control, done, err := router.Run() // run with unimplemented pipe
	r.EqualError(err, "panic")
	r.Nil(control)
	r.Nil(done)
	r.Equal(0, len(router.GetPipes()))
}
