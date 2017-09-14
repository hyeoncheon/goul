package goul_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
)

func Test_Net_10_ClientServerNormal(t *testing.T) {
	r := require.New(t)
	data := "Test Data"

	cmd := make(chan int, 10)
	in := make(chan goul.Item)
	out := make(chan goul.Item)

	var result goul.Item
	packet, err := SetupPacket(data)
	r.NoError(err)

	//** preparing server... ----------------------------------------
	svr, err := goul.NewNetwork("", 6001)
	r.NoError(err)
	svr.SetLogger(goul.NewLogger("debug"))
	go svr.Writer(in)
	time.Sleep(3 * time.Second)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic but recovered")
				return
			}
		}()
		for i := 0; i < 100; i++ {
			in <- packet
			time.Sleep(200 * time.Millisecond)
		}
	}()

	r.NoError(svr.GetError())
	r.Equal(false, svr.InLoop())

	//** preparing client 1... --------------------------------------
	cli, err := goul.NewNetwork("localhost", 6001)
	r.NoError(err)
	cli.SetLogger(goul.NewLogger("debug"))
	go cli.Reader(cmd, out)
	time.Sleep(2 * time.Second)

	r.Equal(true, svr.InLoop())
	r.Equal(true, cli.InLoop())

	for i := 0; i < 5; i++ {
		result = <-out
	}
	cmd <- goul.ComInterrupt

	for ok := true; ok; {
		result, ok = <-out
	}
	time.Sleep(1 * time.Second)

	r.Equal(false, cli.InLoop())
	r.Error(cli.GetError())
	r.Equal(goul.ErrPipeInterrupted, cli.GetError().Error())

	for ok := true; ok; {
		result, ok = <-out
	}

	r.Equal(true, svr.InLoop())
	cli.Close()
	time.Sleep(2 * time.Second)

	r.Equal(false, svr.InLoop())
	r.Error(svr.GetError())
	r.Equal(goul.ErrNetworkConnectionReset, svr.GetError().Error())

	//** preparing client 2... --------------------------------------
	out = make(chan goul.Item)
	cli, err = goul.NewNetwork("localhost", 6001)
	r.NoError(err)
	cli.SetLogger(goul.NewLogger("debug"))
	go cli.Reader(cmd, out)
	time.Sleep(2 * time.Second)

	r.Equal(true, svr.InLoop())
	r.Equal(true, cli.InLoop())

	for i := 0; i < 5; i++ {
		result = <-out
	}

	r.Equal(true, svr.InLoop())
	// instead of cmd <- goul.ComInterrupt
	cli.Close()
	time.Sleep(1 * time.Second)

	r.Equal(false, svr.InLoop())
	r.Error(svr.GetError())
	r.Equal(goul.ErrNetworkConnectionReset, svr.GetError().Error())

	CheckPacket(result, data)

	//close(exitcmd)
	close(in)
	time.Sleep(2 * time.Second)

	r.NotNil(svr.GetError())
	r.Equal(goul.ErrPipeInputClosed, svr.GetError().Error())
	svr.Close()
}

func Test_Net_11_ClientFail(t *testing.T) {
	r := require.New(t)

	cli, err := goul.NewNetwork("localhost", 6006)
	r.Error(err)
	r.Contains(err.Error(), "connection refused")
	r.Nil(cli)
}

func Test_Net_21_ServerRespawn(t *testing.T) {
	r := require.New(t)

	svr1, err := goul.NewNetwork("", 6001)
	r.NoError(err)
	r.NotNil(svr1)

	svr2, err := goul.NewNetwork("", 6001)
	r.Error(err)
	r.Nil(svr2)

	svr1.Close()
}
