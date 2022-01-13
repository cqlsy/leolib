package leofile

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func remove(filename string) bool {
	//os.Remove() // delete record
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// FileExists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// Mkdir
func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// MkdirForFile
func MkdirForFile(path string) (err error) {
	path = filepath.Dir(path)
	return os.MkdirAll(path, os.FileMode(0666))
}

// get file only read
// will create on no exits
func GetFileForRead(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

}

// get file only write
// will create on no exits
// data will append
func GetFileForWrite(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}

// get file only read & write
// will create on no exits
func GetFileForRW(filename string, ) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
}

// Get FileSuffix
func GetSuffix(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

// GetPrefix
func GetPrefix(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[0:i]
		}
	}
	return ""
}

// IsDir
func IsDir(dirname string) bool {
	info, err := os.Stat(dirname)
	return err == nil && info.IsDir()
}

// Copy
func Copy(source string, dest string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	if err == nil {
		si, err := os.Stat(source)
		if err == nil {
			err = os.Chmod(dest, si.Mode())
		}
	}
	return err
}

// DirSize
func DirSize(path string) int64 {
	var dirSize int64 = 0
	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}
	_ = filepath.Walk(path, readSize)
	return dirSize
}

func GetDataFromFile(file io.Reader) (string, error) {
	fd, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(fd), nil
}

func UnZip(dstDirName string, filePath string) {
	// file read
	fr, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()
	// tar read
	tr := tar.NewReader(gr)
	//var subPath = ""
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if HasSuffix(h.Name, "/") {
			//subPath = h.Name
			continue
		}
		var fw *os.File
		if !FileExists(dstDirName + h.Name) {
			err = MkdirForFile(dstDirName + h.Name)
		}
		if err != nil {
			panic(err)
		}
		if FileExists(dstDirName + h.Name) {
			err = os.Remove(dstDirName + h.Name)
		}
		if err != nil {
			panic(err)
		}
		fw, err = os.OpenFile(dstDirName+h.Name, os.O_CREATE|os.O_WRONLY, 0644 /*os.FileMode(h.Mode)*/)
		if err != nil {
			fw, err = os.Open(dstDirName + h.Name)
			panic(err)
		}
		defer fw.Close()
		_, err = io.Copy(fw, tr)
		if err != nil {
			panic(err)
		}
	}
	// delet other
	//d, err := os.ReadDir(dstDirName)
	//subPath = strings.Replace(subPath, "/", "", 1)
	//if err == nil {
	//	for _, f := range d {
	//		if strings.Index(f.Name(), subPath) == -1 {
	//			// delete
	//			fmt.Println("delete", f.Name())
	//			os.RemoveAll(dstDirName + f.Name())
	//		}
	//	}
	//}
}

func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// I think we don't need those function below
//// GetBytes
//func GetBytes(filenameOrURL string, timeout ...time.Duration) ([]byte, error) {
//	if strings.Contains(filenameOrURL, "://") {
//		if strings.Index(filenameOrURL, "leofile://") == 0 {
//			filenameOrURL = filenameOrURL[len("leofile://"):]
//		} else {
//			client := http.DefaultClient
//			if len(timeout) > 0 {
//				client = &http.Client{Timeout: timeout[0]}
//			}
//			r, err := client.Get(filenameOrURL)
//			if err != nil {
//				return nil, err
//			}
//			defer r.Body.Close()
//			if r.StatusCode < 200 || r.StatusCode > 299 {
//				return nil, fmt.Errorf("%d: %s", r.StatusCode, http.StatusText(r.StatusCode))
//			}
//			return ioutil.ReadAll(r.Body)
//		}
//	}
//	return ioutil.ReadFile(filenameOrURL)
//}
//
//// SetBytes
//func SetBytes(filename string, data []byte) error {
//	return ioutil.WriteFile(filename, data, 0666)
//}
//
//// AppendBytes
//func AppendBytes(filename string, data []byte) error {
//	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	_, err = f.Write(data)
//	return err
//}
//
//// GetString
//func GetString(filenameOrURL string, timeout ...time.Duration) (string, error) {
//	bytes, err := GetBytes(filenameOrURL, timeout...)
//	if err != nil {
//		return "", err
//	}
//	return string(bytes), nil
//}
//
//// SetString
//func SetString(filename string, data string) error {
//	return SetBytes(filename, []byte(data))
//}
//
//// AppendString
//func AppendString(filename string, data string) error {
//	return AppendBytes(filename, []byte(data))
//}
//
//// FileTimeModified
//func FileTimeModified(filename string) time.Time {
//	info, err := os.Stat(filename)
//	if err != nil {
//		return time.Time{}
//	}
//	return info.ModTime()
//}
//
//// Find
//func Find(searchDirs []string, filenames ...string) (filePath string, found bool) {
//	for _, dir := range searchDirs {
//		for _, filename := range filenames {
//			filePath = path.Join(dir, filename)
//			if FileExists(filePath) {
//				return filePath, true
//			}
//		}
//	}
//	return "", false
//}
//
