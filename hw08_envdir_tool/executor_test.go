package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	t.Run("Return Zero On Empty Command", func(t *testing.T) {
		code := RunCmd([]string{}, Environment{})
		assert.Equal(t, 0, code)
	})

	t.Run("Return Non-Zero Code On Invalid Command", func(t *testing.T) {
		code := RunCmd([]string{"command-does-not-exists"}, Environment{})
		assert.Equal(t, 1, code)
	})
}
