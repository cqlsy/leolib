package fileUtil

import (
	"errors"
	"github.com/cqlsy/leolib/leofile"
	"github.com/cqlsy/leolib/leotime"
	"github.com/cqlsy/leolib/leoutil"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var filePath = "./"

func Init(rootPath string) {
	filePath = rootPath
}

func SaveFile(src io.Reader, fileName string) (filePath string, err error) {
	filePath = getFilePath(fileName)
	err = leofile.MkdirForFile(filePath)
	if err != nil {
		return
	}
	if leofile.FileExists(filePath) {
		err = os.Remove(filePath)
		if err != nil {
			return filePath, err
		}
	}
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
	// if path start with './', we should delete '.'
	if strings.Index(filePath, "./") == 0 {
		filePath = filePath[1:]
	}
	return
}

func SaveFileWithCallback(src io.Reader, fsize int64, fileName string, fb func(length, downLen int64)) (filePath string, err error) {
	var (
		buf     = make([]byte, 2*1024)
		written int64
	)
	filePath = getFilePath(fileName)
	err = leofile.MkdirForFile(filePath)
	if err != nil {
		return
	}
	if leofile.FileExists(filePath) {
		err = os.Remove(filePath)
		if err != nil {
			return filePath, err
		}
	}
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
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		if fb != nil {
			fb(fsize, written)
		}
	}
	if strings.Index(filePath, "./") == 0 {
		filePath = filePath[1:]
	}
	return
}

func GetFilName(path string) string {
	ss := strings.Split(path, "/")
	if len(ss) <= 0 {
		return "fileName"
	}
	return ss[len(ss)-1]
}

func getFilePath(fileName string) string {
	return getFileFullPath("", fileName)
}

func removePreAndSuf(dir string, old string) string {
	if strings.Index(dir, old) == 0 {
		dir = strings.Replace(dir, old, "", 1)
	}
	if strings.LastIndex(dir, old) == len(dir)-1 {
		dir = dir[:len(dir)-1]
	}
	return dir
}

func getFileFullPath(dir string, fileName string) string {
	fileDir := filePath + "/" + leotime.DateFormat(time.Now(), "YYYY-MM-DD") + "/"
	if dir != "" {
		dir = removePreAndSuf(dir, "/")
		fileDir = filePath + "/" + dir + "/" + leotime.DateFormat(time.Now(), "YYYY-MM-DD") + "/"
	}
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return ""
	}
	name := ""
	i := 0
	houzui := getFileType(fileName)
	for {
		name = leoutil.RandString(16, "default")
		if !leofile.FileExists(fileDir + name + houzui) {
			break
		}
		i++
	}
	return fileDir + name + houzui
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

// example for file download
// fb can set nil
// fileName for only Name
func DownLoadFile(url string, fileName string, fb func(length, downLen int64)) (string, error) {
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	fsize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return "", err
	}
	if resp.Body == nil {
		return "", errors.New("body is null")
	}
	defer resp.Body.Close()
	filePath, err = SaveFileWithCallback(resp.Body, fsize, fileName, fb)
	return filePath, err
}
