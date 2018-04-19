package _internal

import (
	"testing"
)

func Test_Dor(t *testing.T) {
	args := []string{"-l"}
	cmd := newCommand("ls", args, []string{})

	// start
	start(cmd)

	// kill
	kill(cmd)
}
