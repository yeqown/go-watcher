package command

import (
	"testing"
)

func Test_Dor(t *testing.T) {
	cmd := New("ls", []string{"-l"}, []string{})
	// start
	cmd.Start()
	// kill
	cmd.quit()
}
