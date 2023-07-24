package command

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func test()  {
	var (
		rootPath = "/test"
		uploadJsName = "upload.js"
	)
	runCMD("cmd", "/c", "cd", rootPath, "&&", "dir", "/b", "&&", "node", uploadJsName)
}


func runCMD(name string, arg ...string) error {

	//开始执行设定
	cmd := exec.Command(name, arg...)

	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	reader := bufio.NewReader(cmdOut)
	linTemp := ""
	for {
		line, _, err := reader.ReadLine()

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
		linTemp += string(line) + "\n"
	}
	return cmd.Wait()
}
