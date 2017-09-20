package goul_test

import (
	"testing"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
	"github.com/stretchr/testify/require"
)

func Test_BaseRouter_1_Functions(t *testing.T) {
	r := require.New(t)
	var err error
	var router goul.Router
	router = &goul.BaseRouter{}

	err = router.SetLogger(goul.NewLogger("debug"))
	r.NoError(err, "couldn't run SetLogger")

	err = router.SetReader(&goul.BaseAdapter{})
	r.NoError(err, "couldn't run SetReader")

	err = router.SetWriter(&goul.BaseAdapter{})
	r.NoError(err, "couldn't run SetWriter")

	err = router.AddPipe(&goul.BasePipe{}) // BaseRouter does not support
	r.EqualError(err, goul.ErrRouterPipelineNotSupported)

	pipes := router.GetPipes()
	r.Equal(0, len(pipes))

	control, done, err := router.Run()
	r.EqualError(err, goul.ErrAdapterReadNotImplemented)
	r.Nil(control)
	r.Nil(done)
}

// Using DummyAdapter
func Test_BaseRouter_2_ErrorHandling(t *testing.T) {
	r := require.New(t)
	var err error
	var router goul.Router
	router = &goul.BaseRouter{}

	control, done, err := router.Run() // run before set reader and writer
	r.EqualError(err, goul.ErrRouterNoReaderOrWriter)

	err = router.SetReader(&goul.BaseAdapter{})
	control, done, err = router.Run() // run before set writer
	r.EqualError(err, goul.ErrRouterNoReaderOrWriter)

	router = &goul.BaseRouter{}
	err = router.SetWriter(&goul.BaseAdapter{})
	control, done, err = router.Run() // run before set reader
	r.EqualError(err, goul.ErrRouterNoReaderOrWriter)
	r.Nil(control)
	r.Nil(done)

	// Read of BaseAdapter does not implemented.
	err = router.SetReader(&goul.BaseAdapter{})
	control, done, err = router.Run() // run after set reader and writer
	r.EqualError(err, goul.ErrAdapterReadNotImplemented)
	r.Nil(control)
	r.Nil(done)

	// Write of BaseAdapter does not implemented.
	err = router.SetReader(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	control, done, err = router.Run()
	r.EqualError(err, goul.ErrAdapterWriteNotImplemented)
	r.Nil(control)
	r.Nil(done)
}

// Using DummyAdapter
func Test_BaseRouter_3_Run(t *testing.T) {
	r := require.New(t)
	var err error
	var router goul.Router
	router = &goul.BaseRouter{}

	err = router.SetWriter(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	r.NoError(err)
	err = router.SetReader(&adapters.DummyAdapter{Adapter: &goul.BaseAdapter{}})
	r.NoError(err)

	control, done, err := router.Run()
	r.NoError(err)
	r.NotNil(control)
	r.NotNil(done)
	close(control)
	message := <-done
	r.Equal("message", message.String())
}
