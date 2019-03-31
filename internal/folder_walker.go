package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FolderWalker ... to walk all
// not routine-safe for now
type FolderWalker struct {
	// rmDuplicated              bool
	rootDir                   string
	addtionalDirs, resultDirs []string
	excludeChecker            ExcludeChecker
	c                         map[string]bool //
}

// NewFolderWalker ...
func NewFolderWalker(rootDir string, addtionalDirs, excluedDirs []string) *FolderWalker {
	return &FolderWalker{
		// rmDuplicated:  rmDuplicated,
		// excluedDirs:   excluedDirs,
		rootDir:        rootDir,
		addtionalDirs:  addtionalDirs,
		resultDirs:     make([]string, 0),
		excludeChecker: newExcludeChecker(excluedDirs, dirExcludePre),
		c:              make(map[string]bool),
	}
}

func (fw *FolderWalker) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	// not dir or current dir
	if !info.IsDir() || info.Name()[0] == '.' {
		return nil
	}

	// curDir := filepath.Join(path, info.Name())
	curDir := path
	if fw.excludeChecker.Exclude(curDir, dirExcludeJudge) {
		return nil
	}

	// check existed, and if not existed
	if _, ok := fw.c[curDir]; !ok {
		fw.resultDirs = append(fw.resultDirs, curDir)
		fw.c[path] = true
	}

	return nil
}

// Walk ...
func (fw *FolderWalker) Walk() {
	absRoot, err := filepath.Abs(fw.rootDir)
	if err != nil {
		fmt.Println(err)
		return
		// panic(err)
	}

	if err := filepath.Walk(absRoot, fw.walkFunc); err != nil {
		fmt.Println(err)
		return
		// panic(err)
	}

	// walk additional dirs
	if len(fw.addtionalDirs) != 0 {
		for _, dir := range fw.addtionalDirs {
			if dir == "" {
				continue
			}

			if fw.excludeChecker.Exclude(dir, dirExcludeJudge) {
				continue
			}
			if err := filepath.Walk(dir, fw.walkFunc); err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

// Result ... get target dirs
func (fw *FolderWalker) Result() []string {
	return fw.resultDirs
}

// PreFunc ...
type PreFunc func(s string) (string, error)

// JudgeFunc ...
type JudgeFunc func(s string, m map[string]bool) bool

// ExcludeChecker interface to judge an folder in exlude dirs or not
type ExcludeChecker interface {
	Init(excluded []string, pre PreFunc)
	Exclude(s string, jud JudgeFunc) bool
	Excluded() []string
}

// IncludeChecker interface to judge a string is included or not
type IncludeChecker interface {
	Init(included []string, pre PreFunc)
	Include(s string, jud JudgeFunc) bool
	Included() []string
}

var (
	_ ExcludeChecker = &defaultExcludeChecker{}
	_ IncludeChecker = &defaultIncludeChecker{}
)

type defaultExcludeChecker struct {
	cache map[string]bool
}

func newExcludeChecker(excluded []string, pre PreFunc) *defaultExcludeChecker {
	checker := &defaultExcludeChecker{
		cache: make(map[string]bool),
	}

	checker.Init(excluded, pre)

	return checker
}

func (c *defaultExcludeChecker) Init(excluded []string, pre PreFunc) {
	for _, exclude := range excluded {
		predealed, err := pre(exclude)
		if err != nil {
			fmt.Printf("could not prepare with %s\n", exclude)
			continue
		}
		c.cache[predealed] = true
	}
}

func (c *defaultExcludeChecker) Exclude(s string, judge JudgeFunc) bool {
	return judge(s, c.cache)
}

func (c *defaultExcludeChecker) Excluded() []string {
	r := make([]string, len(c.cache))
	cnt := 0
	for k := range c.cache {
		r[cnt] = k
		cnt++
	}
	return r
}

type defaultIncludeChecker struct {
	cache map[string]bool
}

func newIncludeChecker(included []string, pre PreFunc) *defaultIncludeChecker {
	checker := &defaultIncludeChecker{
		cache: make(map[string]bool),
	}

	checker.Init(included, pre)

	return checker
}

func (c *defaultIncludeChecker) Init(included []string, pre PreFunc) {
	for _, include := range included {
		predealed, err := pre(include)
		if err != nil {
			continue
		}
		c.cache[predealed] = true
	}
}

func (c *defaultIncludeChecker) Include(s string, judge JudgeFunc) bool {
	return judge(s, c.cache)
}

// Included ... of defaultIncludeChecker ...
func (c *defaultIncludeChecker) Included() []string {
	r := make([]string, len(c.cache))
	cnt := 0
	for k := range c.cache {
		r[cnt] = k
		cnt++
	}
	return r
}

var (
	_ PreFunc   = dirExcludePre
	_ JudgeFunc = dirExcludeJudge
)

// dirExcludePre ... PreFunc of `dir` default exclude chcker
func dirExcludePre(dir string) (string, error) {
	return filepath.Abs(dir)
}

// dirExcludeJudge ... JudgeFunc of `dir` default exclude chcker
func dirExcludeJudge(dir string, m map[string]bool) bool {
	var (
		absPath string
		err     error
	)
	// get abs path
	if absPath, err = filepath.Abs(dir); err != nil {
		return false
	}

	if _, ok := m[absPath]; ok {
		return true
	}

	// exclude parent dir
	for k := range m {
		if strings.HasPrefix(absPath, k) {
			return true
		}
	}

	return false
}
