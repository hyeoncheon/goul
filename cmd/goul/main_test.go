package main

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hyeoncheon/goul"
)

func Test_Options(t *testing.T) {
	r := require.New(t)
	logger := logger(&Options{})
	r.NotNil(logger)
	_, ok := logger.(goul.Logger)
	r.True(ok)
}
