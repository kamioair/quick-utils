package qservice

import (
	"errors"
	"fmt"
	"github.com/kamioair/quick-utils/qio"
	"runtime"
)

var DeviceCode deviceCode

type deviceCode struct {
}

// LoadFromFile 从文件中获取设备码
func (d *deviceCode) LoadFromFile() (string, error) {
	file := getCodeFile()
	if qio.PathExists(file) {
		code, err := qio.ReadAllString(file)
		if err != nil {
			return "", err
		}
		return code, nil
	}
	return "", errors.New("deviceCode file not find")
}

// SaveToFile 将设备码写入文件
func (d *deviceCode) SaveToFile(code string) error {
	// 写入文件
	file := getCodeFile()
	if file == "" {
		return errors.New("deviceCode file not find")
	}
	err := qio.WriteString(file, code, false)
	if err != nil {
		return err
	}
	return nil
}

func getCodeFile() string {
	root := qio.GetCurrentRoot()
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s\\Program Files\\qf\\device", root)
	case "linux":
		return "/dev/qf/device"
	}
	return ""
}
