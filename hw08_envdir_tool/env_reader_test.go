package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	mainDir := "/tmp/"
	dir, _ := os.MkdirTemp(mainDir, "env_reader")
	defer os.RemoveAll(dir)

	t.Run("incorrect file name", func(t *testing.T) {
		file, _ := os.CreateTemp(dir, "tes=t.")
		defer os.Remove(file.Name())
		_, _ = file.WriteString("incorrect name file")

		dirEnv, err := ReadDir(dir)

		require.Equal(t, Environment{}, dirEnv)
		require.Error(t, err)
	})

	t.Run("check valid value", func(t *testing.T) {
		file, _ := os.CreateTemp(dir, "FOO.")
		defer os.Remove(file.Name())
		_, _ = file.WriteString("test value")
		fileInfo, _ := file.Stat()

		dirEnv, err := ReadDir(dir)

		require.Equal(t, Environment{
			fileInfo.Name(): EnvValue{
				Value:      "test value",
				NeedRemove: false,
			},
		}, dirEnv)
		require.NoError(t, err)
	})

	t.Run("check clear value", func(t *testing.T) {
		file, _ := os.CreateTemp(dir, "FOO.")
		defer os.Remove(file.Name())
		_, _ = file.WriteString("test value  \t   \n new line text")
		fileInfo, _ := file.Stat()

		dirEnv, err := ReadDir(dir)

		require.Equal(t, Environment{
			fileInfo.Name(): EnvValue{
				Value:      "test value",
				NeedRemove: false,
			},
		}, dirEnv)
		require.NoError(t, err)
	})
}
