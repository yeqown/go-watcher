package utils

/*
 * util functions
 */

import (
	"os"
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

// // CheckFileIncluded 根据后缀检测
// func CheckFileIncluded(filename string, ftyps []string) bool {
// 	for _, wft := range ftyps {
// 		if strings.HasSuffix(filename, wft) {
// 			return true
// 		}
// 	}
// 	return false
// }

// // CheckFileRegexpExcluded 检查文件名字是否匹配正则
// func CheckFileRegexpExcluded(filename string, regexps []string) bool {
// 	for _, reg := range regexps {
// 		if m, err := regexp.Match(reg, []byte(filename)); err != nil {
// 			// panic(err)
// 			continue
// 		} else {
// 			// cond: err == nil
// 			if m {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }
