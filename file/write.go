package file

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 替换一整行
func  ReplaceFile(filePath, origin, target string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	data := make([]string, 0)
	for {
		//读取每一行内容
		tmp, _, err := reader.ReadLine()

		if err != nil {
			//读到末尾
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return err
			}
		}
		line := string(tmp) + "\n"

			if strings.Contains(line, origin) {
				line = target
			}


		data = append(data, line)
	}

	f.Close()
	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return WriteStringToFile(filePath, strings.Join(data, ""))
}

// WriteStringToFile 写文件
func  WriteStringToFile(path, data string) (err error) {
	//1.判断文件是否存在
	has := fileExists(path)
	if !has {
		_, err = creatNestedFile(path)
		if err != nil {
			return
		}
	}

	//2、打开文件
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	//3、写文件
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

// base64ToFile 把base64写入文件中
func  base64ToFile(filePath, data string) (err error) {
	decodeData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()
	n, err := f.Write(decodeData)
	if err == nil && n != len(decodeData) {
		return errors.New("没有写成功")
	}
	return
}

// 解压文件
func unzipFile(filePath string) (err error){
	 var dst      = "apk/"
	//解压文件
	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return
	}

	var decodeName string
	for _, f := range archive.File {
		if f.Flags == 0 {
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			//如果标志为是 1 << 11也就是 2048  则是utf-8编码
			decodeName = f.Name
		}
		filePath := filepath.Join(dst, decodeName)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {

		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	archive.Close()
	return os.RemoveAll(filePath)
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}