package ioutools

import "os"

// CheckDirOrFileIsExist 检查目录或者文件是否存在
// @auth      iouAkira
// @param     path string 文件夹或者文件的绝对路径
func CheckDirOrFileIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
