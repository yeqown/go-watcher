package internal

import (
	"testing"
)

func Test_NewWatcher(t *testing.T) {
	exit := make(chan bool)
	watchingFiletypes := []string{"go"}
	unwatchingRegular := []string{}

	if w, err := NewWatcher([]string{"."}, exit,
		watchingFiletypes, unwatchingRegular); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		w.Exit()
	}
}
