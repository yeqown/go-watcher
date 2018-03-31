// // package gowatch
package internal

import (
	"testing"
)

func Test_fileExist(t *testing.T) {
	fpath := "./testdata/testdata_inner/not_ex_file"
	fpath_ex := "./testdata/test_gowatch.txt"

	if fileExist(fpath) {
		t.Logf("%s shouldn't be ex\n", fpath)
		t.Fail()
	}

	if !fileExist(fpath_ex) {
		t.Logf("%s should be ex\n", fpath)
		t.Fail()
	}
}

func Test_isFileExclued(t *testing.T) {
	// excluded
	fname := "./testdata/test_gowatch.go"
	fname2 := "./testdata/test_gowatch.txt"
	fname3 := "./testdata/testdata_inner/test_gowatch.go"
	// not in
	fname4 := "./dor.go"

	excluedPaths := []string{
		"./testdata",
		"./testdata/testdata_inner",
	}
	if !checkPathExcluded(fname, excluedPaths) {
		t.Logf("%s is in excludePaths\n", fname)
		t.Fail()
	}

	if !checkPathExcluded(fname2, excluedPaths) {
		t.Logf("%s is in excludePaths\n", fname2)
		t.Fail()
	}

	if !checkPathExcluded(fname3, excluedPaths) {
		t.Logf("%s is in excludePaths\n", fname3)
		t.Fail()
	}

	if checkPathExcluded(fname4, excluedPaths) {
		t.Logf("%s is not in excludePaths\n", fname4)
		t.Fail()
	}
}

func Test_checkFileRegexpExcluded(t *testing.T) {
	exclude_regexps := []string{
		".go$",
	}
	fname := "text_gowatch.txt"
	fname_exclu := "test_gowatch.go"

	if !checkFileRegexpExcluded(fname_exclu, exclude_regexps) {
		t.Logf("%s is exclude regexp\n", fname_exclu)
		t.Fail()
	}

	if checkFileRegexpExcluded(fname, exclude_regexps) {
		t.Logf("%s is not exclude regexp\n", fname)
		t.Fail()
	}
}

func Test_checkDirIncluded(t *testing.T) {

	AppendWatchFiletypes(".txt")

	if !checkFileIncluded("./testdata/test_gowatch.go", WatchFiletypes) {
		t.Log("./testdata/test_gowatch.txt is included")
		t.Fail()
	}
	if !checkFileIncluded("./testdata/test_gowatch.txt", WatchFiletypes) {
		t.Log("./testdata/test_gowatch.txt is included")
		t.Fail()
	}

	if checkFileIncluded("./config.yml", WatchFiletypes) {
		t.Log("./config.yml is not watched")
		t.Fail()
	}
}

func Test_WalkDirectoryRecursive(t *testing.T) {
	excluedPaths := []string{
		"./testdata/testdata_inner",
	}
	paths := []string{}
	WalkDirectoryRecursive("./testdata", excluedPaths, &paths)
	for _, p := range paths {
		t.Log("path", p)
	}
}
