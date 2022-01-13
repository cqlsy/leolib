package leofile

import (
	"testing"
)

func Test_openFile(t *testing.T) {
	file, err := GetFileForRead("../data/test.txt")
	if err == nil {
		println(file.Name())
	} else {
		println(err.Error())
	}
}
