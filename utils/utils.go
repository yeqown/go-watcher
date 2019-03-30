package utils

/*
 * util functions
 */

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// IsFileExist check file exist
func IsFileExist(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil || os.IsExist(err)
}

// GetFileModTime ... get file modified unix time
func GetFileModTime(path string) int64 {

	path = strings.Replace(path, "\\", "/", -1)

	if f, err := os.Open(path); err == nil {
		defer f.Close()
		if fi, err := f.Stat(); err == nil {
			return fi.ModTime().Unix()
		}
	}
	return time.Now().Unix()
}

// WalkDirectory ...
//
// walk directory to get all paths to directory and sub-directory recursive
//
// TODO: remove duplicated path in paths
// TODO: support callback function
func WalkDirectory(dir string, excluedPaths []string, paths *[]string, recursive bool) {
	var (
		fInfos []os.FileInfo
		err    error
	)

	if fInfos, err = ioutil.ReadDir(dir); err != nil {
		return
	}

	for _, finfo := range fInfos {
		tmpDir := filepath.Join(dir, finfo.Name())

		if finfo.IsDir() && checkPathExcluded(tmpDir, excluedPaths) {
			continue
		}
		if finfo.IsDir() && finfo.Name()[0] != '.' && recursive {
			WalkDirectory(tmpDir, excluedPaths, paths, recursive)
			*paths = append(*paths, tmpDir)
		}
	}

}

// UnixTimeDuration 计算时间戳的间隔
func UnixTimeDuration(t1, t2 int64) time.Duration {
	d := t1 - t2
	if d < 0 {
		return time.Duration(-d)
	}
	return time.Duration(d)

	// tt1 := time.Unix(t1, 0)
	// tt2 := time.Unix(t2, 0)
	// if t1 > t2 {
	// 	return tt1.Sub(tt2)
	// }
	// return tt2.Sub(tt1)
}

// 检查文件夹是否需要排除
func checkPathExcluded(fpath string, excluedPaths []string) bool {
	var (
		err      error
		absPath  string
		absFPath string
	)
	// 获取绝对路径
	if absFPath, err = filepath.Abs(fpath); err != nil {
		return false
	}

	for _, exclPath := range excluedPaths {
		// 文件夹名字 和 排除名字一致
		if fpath == exclPath || absFPath == exclPath {
			return true
		}
		if absPath, err = filepath.Abs(exclPath); err != nil {
			continue
		}
		// println(absPath, absath)
		if strings.HasPrefix(absFPath, absPath) {
			return true
		}
	}
	return false
}

// CheckFileIncluded 根据后缀检测
func CheckFileIncluded(filename string, ftyps []string) bool {
	for _, wft := range ftyps {
		if strings.HasSuffix(filename, wft) {
			return true
		}
	}
	return false
}

// CheckFileRegexpExcluded 检查文件名字是否匹配正则
func CheckFileRegexpExcluded(filename string, regexps []string) bool {
	for _, reg := range regexps {
		if m, err := regexp.Match(reg, []byte(filename)); err != nil {
			// panic(err)
			continue
		} else {
			// cond: err == nil
			if m {
				return true
			}
		}
	}
	return false
}
