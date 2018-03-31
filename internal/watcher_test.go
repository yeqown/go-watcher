// package gowatch
package internal

import (
	"testing"
)

func Test_NewWatcher(t *testing.T) {
	_, err := NewWatcher()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
