package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envName := "NEW"
	cmd := []string{"echo", ""}

	t.Run("check change env", func(t *testing.T) {
		os.Setenv(envName, "old")

		returnCode := RunCmd(cmd, Environment{
			envName: EnvValue{
				Value:      "new",
				NeedRemove: false,
			},
		})

		require.Equal(t, "new", os.Getenv(envName))
		require.Equal(t, 0, returnCode)
	})

	t.Run("check del env", func(t *testing.T) {
		os.Setenv(envName, "old")

		returnCode := RunCmd(cmd, Environment{
			envName: EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		})

		_, ok := os.LookupEnv(envName)

		require.False(t, ok)
		require.Equal(t, 0, returnCode)
	})

	t.Run("check create env", func(t *testing.T) {
		returnCode := RunCmd(cmd, Environment{
			envName: EnvValue{
				Value:      "old",
				NeedRemove: false,
			},
		})

		_, ok := os.LookupEnv(envName)

		require.True(t, ok)
		require.Equal(t, 0, returnCode)
	})
}
