package adapters_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/adapters"
	. "github.com/hyeoncheon/goul/testing"
)

func Test_Network_10_Normal(t *testing.T) {
	r := require.New(t)

	control0, outServer := debugServer(r)
	r.NotNil(control0)
	control1, done1 := generatorClient(r, "C1")
	control2, done2 := generatorClient(r, "C2")

	time.Sleep(1000 * time.Millisecond)
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		control1 <- &goul.ItemGeneric{Meta: "packet", DATA: []byte("TD1")}
		<-outServer
		//! replace with compressed data case
		time.Sleep(100 * time.Millisecond)
		control2 <- &goul.ItemGeneric{Meta: "packet", DATA: []byte("TD2")}
		out := <-outServer
		// check data integrity
		p, ok := out.(gopacket.Packet)
		r.True(ok)
		r.NotNil(p)
		r.NotNil(p.ApplicationLayer())
		r.Equal("TD2", string(p.ApplicationLayer().Payload()))
	}
	close(control1)
	<-done1 //! check status of client
	close(control2)
	<-done2 //! check status of client
	time.Sleep(1000 * time.Millisecond)
	close(control0)
	<-outServer //! check status of server
}

func Test_Network_20_Interrupted(t *testing.T) {
	r := require.New(t)

	control0, outServer := debugServer(r)
	r.NotNil(control0)
	control1, done1 := generatorClient(r, "C1")

	time.Sleep(1000 * time.Millisecond)
	close(control0)
	<-outServer

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		<-done1
		close(control1) //! unusal. anyway close control channel
		wg.Done()
		return
	}()

	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		control1 <- &goul.ItemGeneric{Meta: "packet", DATA: []byte("TD1")}
		<-outServer
	}

	wg.Wait()
	time.Sleep(1000 * time.Millisecond)
}

func Test_Network_21_Exceptions(t *testing.T) {
	r := require.New(t)

	reader, err := adapters.NewNetwork("", 600)
	r.NoError(err)
	server := &goul.BaseRouter{}
	server.SetReader(reader)
	server.SetWriter(&GeneratorAdapter{ID: "  --SW", Adapter: &goul.BaseAdapter{}})
	_, _, err = server.Run()
	r.Error(err)
	r.Contains(err.Error(), "permission denied")
	err = reader.Close()
	r.NoError(err)

	writer, err := adapters.NewNetwork("localhost", 600)
	r.NoError(err)
	client := &goul.BaseRouter{}
	client.SetReader(&GeneratorAdapter{ID: "    ", Adapter: &goul.BaseAdapter{}})
	client.SetWriter(writer)
	_, _, err = client.Run()
	r.Error(err)
	r.Contains(err.Error(), "connection refused")
	err = writer.Close()
	r.NoError(err)
}

func Test_Network_22_Close(t *testing.T) {
	r := require.New(t)

	reader, err := adapters.NewNetwork("", 6006)
	r.NoError(err)
	server := &goul.BaseRouter{}
	server.SetReader(reader)
	server.SetWriter(&GeneratorAdapter{ID: "  --SW", Adapter: &goul.BaseAdapter{}})
	control0, outServer, err := server.Run()
	r.NoError(err)

	close(control0)
	<-outServer
	reader.Close()
}

func Test_Network_23_Respawn(t *testing.T) {
	r := require.New(t)

	c, d := debugServer(r)

	reader, err := adapters.NewNetwork("", 6006)
	r.NoError(err)
	server := &goul.BaseRouter{}
	server.SetReader(reader)
	server.SetWriter(&GeneratorAdapter{ID: "  --SW", Adapter: &goul.BaseAdapter{}})
	_, _, err = server.Run()
	r.Error(err)
	r.Contains(err.Error(), "already in use")

	close(c)
	<-d
}

//** utilities

func debugServer(r *require.Assertions) (control, out chan goul.Item) {
	reader, err := adapters.NewNetwork("", 6006)
	reader.ID = "  ->SR"
	r.NoError(err)
	server := &goul.BaseRouter{}
	server.SetLogger(goul.NewLogger("debug"))
	server.SetReader(reader)
	server.SetWriter(&GeneratorAdapter{ID: "  --SW", Adapter: &goul.BaseAdapter{}})
	control, out, err = server.Run()
	r.NoError(err)

	return control, out
}

func generatorClient(r *require.Assertions, name string) (control, out chan goul.Item) {
	writer, err := adapters.NewNetwork("localhost", 6006)
	writer.ID = name + "->  "
	r.NoError(err)
	client := &goul.BaseRouter{}
	client.SetLogger(goul.NewLogger("debug"))
	client.SetReader(&GeneratorAdapter{ID: name + "    ", Adapter: &goul.BaseAdapter{}})
	client.SetWriter(writer)
	control, done, err := client.Run()
	r.NoError(err)

	return control, done
}
