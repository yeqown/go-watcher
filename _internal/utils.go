/*
 * util functions
 */

package _internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	WatchFiletypes = []string{".go"} // 遍历文件夹的时候需要忽略的文件夹名字
	UnWatchRegExps = []string{}
	emptyPaths     = []string{} // 默认空路径列表
)

// check file exist
func fileExist(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil || os.IsExist(err)
}

func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
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

// 根据后缀检测
func checkFileIncluded(name string, ftyps []string) bool {
	for _, wft := range ftyps {
		if strings.HasSuffix(name, wft) {
			return true
		}
	}
	return false
}

// 检查文件名字是否匹配正则
func checkFileRegexpExcluded(filename string, regexps []string) bool {
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

// 增加监听的文件类型
func AppendWatchFiletypes(names ...string) {
	WatchFiletypes = append(WatchFiletypes, names...)
}

// 增加不需要监视的正则匹配表达式
func AppendUnWatchRegexps(exps ...string) {
	UnWatchRegExps = append(UnWatchRegExps, exps...)
}

/*
 * 递归遍历文件夹，返回所有路径，不包括需要排除的文件夹及子文件夹
 * ? 返回的内容里面有重复的路径？
 */
func WalkDirectoryRecursive(dir string, excluedPaths []string, paths *[]string) {
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
		if finfo.IsDir() && finfo.Name()[0] != '.' {
			WalkDirectoryRecursive(tmpDir, excluedPaths, paths)
			*paths = append(*paths, tmpDir)
		}
	}

}

// 计算时间戳的间隔
func UnixTimeDuration(t1, t2 int64) time.Duration {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)
	if t1 > t2 {
		return tt1.Sub(tt2)
	}
	return tt2.Sub(tt1)
}
