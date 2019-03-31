package internal

import (
	"testing"
)

func TestDefaultExcludeChecker(t *testing.T) {
	excluded := []string{
		"exclude",
		"/home/testdata",
	}
	checker := newExcludeChecker(excluded, dirExcludePre)

	if !checker.Exclude("./exclude", dirExcludeJudge) {
		t.Error("wrong judge './exclude'")
		t.FailNow()
	}

	if !checker.Exclude("/home/testdata/demo", dirExcludeJudge) {
		t.Error("wrong judge '/home/testdata/demo'")
		t.FailNow()
	}

	if !checker.Exclude("./exclude/demo", dirExcludeJudge) {
		t.Error("wrong judge './exclude/demo'")
		t.FailNow()
	}

	if checker.Exclude("./include", dirExcludeJudge) {
		t.Error("wrong judge of './include'")
		t.FailNow()
	}

}
func TestFolderWalker(t *testing.T) {
	excluded := []string{
		"./exclude",
		"/home/testdata",
	}
	additional := []string{""}

	walker := NewFolderWalker("", additional, excluded)
	walker.Walk()
	paths := walker.Result()

	t.Log(paths)
}
