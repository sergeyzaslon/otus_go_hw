package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var source = "testdata/env"

func TestReadDir(t *testing.T) {
	t.Run("Check path", func(t *testing.T) {
		env, err := ReadDir(source)
		require.Nil(t, err)

		expected := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}

		require.Equal(t, expected, env)
	})
}
