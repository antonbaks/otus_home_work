package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	pathFileTo := "out.txt"

	c := Copy("/dev/urandom", pathFileTo, 0, 0)

	require.Truef(t, errors.Is(c, ErrUnsupportedFile), "actual error %q", c)

	c = Copy("testdata/input.txt", pathFileTo, -1, 0)

	require.Truef(t, errors.Is(c, ErrSeek), "actual error %q", c)

	c = Copy("testdata/input.txt", pathFileTo, 10000, 0)

	require.Truef(t, errors.Is(c, ErrOffsetExceedsFileSize), "actual error %q", c)

	c = Copy("testdata/input.txt", pathFileTo, 0, 99999999)

	require.Equal(t, c, nil)

	_ = os.Remove(pathFileTo)
}
