package goul_test

import (
	"testing"

	"github.com/hyeoncheon/goul"
	"github.com/stretchr/testify/require"
)

func Test_PrintDevices_1_Normal(t *testing.T) {
	r := require.New(t)

	err := goul.PrintDevices()
	r.NoError(err)
}
