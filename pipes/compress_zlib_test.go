package pipes_test

import (
	"testing"
	"time"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
	"github.com/stretchr/testify/require"
)

func Test_CompressZLib(t *testing.T) {
	pt := &PacketPipeTest{PacketPipe: &pipes.CompressZLib{}, T: t}
	pt.Run()
}

func Test_CompressZLib_MalformData(t *testing.T) {
	pipe := &pipes.CompressZLib{}
	in := make(chan goul.Item)
	out := make(chan goul.Item)

	go pipe.Reverse(in, out)

	in <- &goul.ItemGeneric{Meta: "application/zlib", DATA: []byte{1}}

	for i := 0; i < 20 && pipe.GetError() == nil; i++ {
		time.Sleep(100 * time.Millisecond)
	}
	require.New(t).Equal("CouldNotCreateNewReader", pipe.GetError().Error())
	close(in)
}
