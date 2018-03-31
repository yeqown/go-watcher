// package gowatch
package internal

import (
	"testing"
)

func Test_Dor(t *testing.T) {
	args := []string{
		"-l",
	}
	cmd := newCommand("ls", args, EmptyEnvs)

	// start
	start(cmd)

	// kill
	kill(cmd)
}
