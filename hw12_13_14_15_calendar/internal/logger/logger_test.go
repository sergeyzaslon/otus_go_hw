package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("test dev", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(file.Name(), "debug", "text_simple")
		l.Debug("DEBUG %s", "one")
		l.Info("INFO %s", "two")
		l.Warn("WARNING %s", "three")
		l.Error("ERROR %s", "four")

		logContent, _ := os.ReadFile(file.Name())
		logExpected := "debug\tDEBUG one\ninfo\tINFO two\nwarning\tWARNING three\nerror\tERROR four\n"

		require.Equal(t, logExpected, string(logContent))
	})

	t.Run("test prod", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(file.Name(), "warn", "text_simple")
		l.Debug("DEBUG %s", "one")
		l.Info("INFO %s", "two")
		l.Warn("WARNING %s", "three")
		l.Error("ERROR %s", "four")

		// file.Close()
		logContent, _ := os.ReadFile(file.Name())
		fmt.Println(string(logContent))
		logExpected := "warning\tWARNING three\nerror\tERROR four\n"

		require.Equal(t, logExpected, string(logContent))
	})
}
