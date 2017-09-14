package goul_test

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
)

func Test_PrintDevices_1_Normal(t *testing.T) {
	user, err := user.Current()
	if err != nil {
		return
	}
	isTravis := user.Username == "travis"

	r := require.New(t)

	err = goul.PrintDevices()
	if isTravis {
		r.Error(err)
	} else {
		r.NoError(err)
	}
}
