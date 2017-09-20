package pipes_test

import (
	"testing"

	"github.com/hyeoncheon/goul"
	"github.com/hyeoncheon/goul/pipes"
	"github.com/stretchr/testify/require"
)

func Test_CompressZLib(t *testing.T) {
	pts := &PipeTestSuite{
		C: &pipes.CompressZLib{Pipe: &goul.BasePipe{Mode: goul.ModeConverter}},
		R: &pipes.CompressZLib{Pipe: &goul.BasePipe{Mode: goul.ModeReverter}},
		T: t,
	}
	pts.Run()

	ptsda := &PipeTestSuiteDirectAccess{
		C: &pipes.CompressZLib{},
		R: &pipes.CompressZLib{},
		T: t,
	}
	ptsda.Run()
}

func Test_CompressZLib_MalformData(t *testing.T) {
	r := require.New(t)
	pipe := &pipes.CompressZLib{Pipe: &goul.BasePipe{Mode: goul.ModeConverter}}

	in := make(chan goul.Item)
	out, err := pipe.Revert(in, nil)
	r.NoError(err)

	in <- &goul.ItemGeneric{Meta: "application/zlib", DATA: []byte{1}}
	// this will not generate any output.
	in <- &goul.ItemGeneric{Meta: "dummy data", DATA: []byte{1}}
	<-out // it means, previous processing also done. (hack for timing)
	r.Equal(pipes.CouldNotCreateNewZLibReader, pipe.GetError().Error())

	close(in)
	<-out
	r.Equal(goul.ErrPipeInputClosed, pipe.GetError().Error())
}
