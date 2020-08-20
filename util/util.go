package util

import "os"

// 获取项目路径
func LocalPath() string {
	path, _ := os.Getwd()
	return path+"/../"
}