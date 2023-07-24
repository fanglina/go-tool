package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// 判断所给路径文件/文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		fmt.Println(err)
		return false
	}
	return true
}


// CreatNestedFile 给定path创建文件，如果目录不存在就递归创建
func creatNestedFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if !fileExists(basePath) {
		err := os.MkdirAll(basePath, 0700)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(path)
}
