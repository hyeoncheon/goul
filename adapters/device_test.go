package adapters_test

import (
	"testing"

	"github.com/hyeoncheon/goul"

	"github.com/hyeoncheon/goul/adapters"

	"github.com/stretchr/testify/require"
)

func Test_DeviceAdapter_1_NormalFlow(t *testing.T) {
	var err error
	var adapter goul.Adapter
	r := require.New(t)

	adapter, err = adapters.NewDevice("eth0")
	r.NoError(err)
	r.NotNil(adapter)

	logger := goul.NewLogger("debug")
	err = adapter.SetLogger(logger)
	r.NoError(err)

	lgr := adapter.GetLogger()
	r.NotNil(lgr)

	in := make(chan goul.Item)
	go func() {
		in <- &goul.ItemGeneric{Meta: "Message", DATA: []byte{1}}
		return
	}()

	out, err := adapter.Read(in, nil)
	r.NoError(err)
	<-out
	r.Contains(adapter.GetError().Error(), "Permission Denied")

	out, err = adapter.Write(in, nil)
	r.NoError(err)
	<-out
	r.Contains(adapter.GetError().Error(), "Permission Denied")

	adapter.Close()
}

func Test_DeviceAdapter_2_AdapterSpecific(t *testing.T) {
	r := require.New(t)

	// get instance as adapters.DeviceAdapter instead of goul.Adapter
	adapter, err := adapters.NewDevice("eth0")
	r.NoError(err)
	r.NotNil(adapter)

	err = adapter.SetFilter("port 80")
	r.NoError(err)
	err = adapter.SetOptions(false, 1500, 1)
	r.NoError(err)
}

func Test_DeviceAdapter_3_Uninitialized(t *testing.T) {
	r := require.New(t)
	var err error

	//! for this case, did not inherit base adapter, it cannot be used
	//! properly. need to be enhanced this limitation.
	adapter := &adapters.DeviceAdapter{}
	err = adapter.SetOptions(false, 1500, 1)
	r.EqualError(err, adapters.ErrDeviceAdapterNotInitialized)

	in := make(chan goul.Item)
	out, err := adapter.Read(in, nil)
	<-out

	//! using DeviceAdapter with BaseAdapter but not initialize.
	adapter = &adapters.DeviceAdapter{Adapter: &goul.BaseAdapter{}}
	err = adapter.SetOptions(false, 1500, 1)
	r.EqualError(err, adapters.ErrDeviceAdapterNotInitialized)

	out, err = adapter.Read(in, nil)
	<-out
	r.EqualError(adapter.GetError(), adapters.ErrDeviceAdapterNotInitialized)
}
