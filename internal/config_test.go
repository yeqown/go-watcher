// package gowatch
package internal

import (
	"testing"
)

func Test_parseConfig(t *testing.T) {
	ParseConfig("./testdata/gowatch.yml")
	t.Logf("%v", cfg)
}

func Test_outputConf(t *testing.T) {
	if err := OutputDefaultConf("./testdata/gowatch.yml"); err != nil {
		t.Log(err)
		t.Fail()
	}
}
