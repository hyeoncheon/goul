package goul_test

import (
	"testing"
	"time"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
	"github.com/stretchr/testify/require"
)

func Test_Goul_10_Sender(t *testing.T) {
	r := require.New(t)
	//data := "Test Data"

	//***
	cmd := make(chan int)
	//in := make(chan goul.Item)
	//out := make(chan goul.Item)

	gl, err := goul.New("eth0", false, cmd, true)
	r.NoError(err)
	r.NotNil(gl)
	err = gl.SetOptions(false, 1600, 1)
	r.NoError(err)
	err = gl.SetLogger(goul.NewLogger("debug"))
	r.NoError(err)
	err = gl.SetFilter("port 80")
	r.NoError(err)
	err = gl.SetReader(gl)
	r.NoError(err)

	gl.AddPipe(&pipes.PacketPrinter{})
	gl.AddPipe(&pipes.DataCounter{})

	r.NoError(gl.GetError())
	r.Equal(false, gl.InLoop())

	// FAIL, without writer
	err = gl.Run()
	r.Error(err)
	r.Contains(err.Error(), "no reader or writer")

	r.Error(gl.GetError())
	r.Equal(goul.ErrNoReaderOrWriter, gl.GetError().Error())
	r.Equal(false, gl.InLoop())

	gl.SetWriter(&pipes.NullWriter{})
	err = gl.Run()
	r.NoError(err)
	//r.Contains(err.Error(), "permissiii")

	time.Sleep(1 * time.Second)
	//cmd <- goul.ComInterrupt
	/*

		gl.Writer(in)

		in <- setupPacket(r, data)
		close(in)

		//	cli, err := goul.New("localhost", 6001)
		//	r.NoError(err)
		//	r.NotNil(cli)

		//go svr.Writer(in)
		//go cli.Reader(cmd, out)
		/*

			result := <-out
			utilCheckPacket(r, result, data)

			cmd <- goul.ComInterrupt
	*/
	gl.Close()
}

func Test_Goul_20_Receiver(t *testing.T) {
	r := require.New(t)
	//data := "Test Data"

	//***
	cmd := make(chan int)
	//in := make(chan goul.Item)
	//out := make(chan goul.Item)

	gl, err := goul.New("eth0", true, cmd, true)
	r.NoError(err)
	r.NotNil(gl)
	err = gl.SetOptions(false, 1600, 1)
	r.NoError(err)
	err = gl.SetLogger(goul.NewLogger("debug"))
	r.NoError(err)

	gl.AddPipe(&pipes.PacketPrinter{})
	gl.AddPipe(&pipes.DataCounter{})

	r.NoError(gl.GetError())
	r.Equal(false, gl.InLoop())

	// FAIL, without reader
	err = gl.Run()
	r.Error(err)
	r.Contains(err.Error(), "no reader or writer")

	r.Error(gl.GetError())
	r.Equal(goul.ErrNoReaderOrWriter, gl.GetError().Error())
	r.Equal(false, gl.InLoop())

	err = gl.SetReader(gl)
	r.NoError(err)
	err = gl.Run()
	r.NoError(err)
	//r.Contains(err.Error(), "permissiii")

	time.Sleep(1 * time.Second)
	//cmd <- goul.ComInterrupt
	/*

		gl.Writer(in)

		in <- setupPacket(r, data)
		close(in)

		//	cli, err := goul.New("localhost", 6001)
		//	r.NoError(err)
		//	r.NotNil(cli)

		//go svr.Writer(in)
		//go cli.Reader(cmd, out)
		/*

			result := <-out
			utilCheckPacket(r, result, data)

			cmd <- goul.ComInterrupt
	*/
	gl.Close()
}
