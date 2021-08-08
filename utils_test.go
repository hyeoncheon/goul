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
		// Ah... I don't remember but it seems like this was an workaround
		// for Travis CI. but not it works as the same as local env.
		// just leave it as is as a history.
		//r.Error(err)
		r.NoError(err)
	} else {
		r.NoError(err)
	}
}
