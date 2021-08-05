package yeefile

import (
	"io"
	"lib/yeeUtil"
	"lib/yeetime"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

var filePath = "./data"

// 上传的文件存放在哪里 ，默认是当前文件加下面，所以设置的时候，不能设置 ./
func Init(rootPath string) {
	filePath = rootPath
}

func SaveFile(fileHeader *multipart.FileHeader) (filePath string, err error) {
	src, err := fileHeader.Open()
	if err != nil {
		return
	}
	defer src.Close()
	filePath = getFilePath(fileHeader.Filename)
	// Destination
	dst, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return
	}
	return
}

func getFilePath(fileName string) string {
	fileDir := filePath + "/" + yeetime.DateFormat(time.Now(), "YYYY-MM-DD") + "/"
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return ""
	}
	name := ""
	i := 0
	houzui := getFileType(fileName)
	for {
		name = yeeUtil.RandString(16, "default")
		if !FileExists(fileDir + name + houzui) {
			break
		}
		i++
	}
	return fileDir + name + houzui
}

func GetFilName(path string) string {
	ss := strings.Split(path, "/")
	if len(ss) <= 0 {
		return "fileName"
	}
	return ss[len(ss)-1]
}

func getFileType(fileName string) string {
	re := ""
	if fileName != "" {
		index := strings.LastIndex(fileName, ".")
		if index != -1 {
			re = fileName[index:]
		}
	}
	return re
}
