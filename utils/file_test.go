package utils

import (
	"testing"
)

func Test_IsFileExist(t *testing.T) {
	fpath := "./testdata/testdata_inner/not_ex_file"
	fpathEx := "./testdata/test_gowatch.txt"

	if IsFileExist(fpath) {
		t.Logf("%s shouldn't be ex\n", fpath)
		t.Fail()
	}

	if !IsFileExist(fpathEx) {
		t.Logf("%s should be ex\n", fpath)
		t.Fail()
	}
}

// func Test_CheckFileRegexpExcluded(t *testing.T) {
// 	excludeRegexps := []string{
// 		".go$",
// 	}
// 	fname := "text_gowatch.txt"
// 	fnameExclu := "test_gowatch.go"

// 	if !CheckFileRegexpExcluded(fnameExclu, excludeRegexps) {
// 		t.Logf("%s is exclude regexp\n", fnameExclu)
// 		t.Fail()
// 	}

// 	if CheckFileRegexpExcluded(fname, excludeRegexps) {
// 		t.Logf("%s is not exclude regexp\n", fname)
// 		t.Fail()
// 	}
// }

// func Test_CheckDirIncluded(t *testing.T) {

// 	watchFiletypes := []string{"go", "txt"}

// 	if !CheckFileIncluded("./testdata/test_gowatch.go", watchFiletypes) {
// 		t.Log("./testdata/test_gowatch.txt is included")
// 		t.Fail()
// 	}
// 	if !CheckFileIncluded("./testdata/test_gowatch.txt", watchFiletypes) {
// 		t.Log("./testdata/test_gowatch.txt is included")
// 		t.Fail()
// 	}

// 	if CheckFileIncluded("./config.yml", watchFiletypes) {
// 		t.Log("./config.yml is not watched")
// 		t.Fail()
// 	}
// }

// func Test_WalkDirectory(t *testing.T) {
// 	excluedPaths := []string{
// 		"./testdata/testdata_inner",
// 	}
// 	paths := []string{}
// 	WalkDirectory("./testdata", excluedPaths, &paths, true)
// 	for _, p := range paths {
// 		t.Log("path", p)
// 	}
// }
