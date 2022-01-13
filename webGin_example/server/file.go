package server

import (
	"github.com/gin-gonic/gin"
	"github.com/cqlsy/leolib/leowebgin"
	"github.com/cqlsy/leolib/webGin_example/server/base"
)

type File struct {
}

func (File) UploadFile(c *gin.Context) {
	filePath, err := leowebgin.UploadFile(c, "file")
	if err != nil {
		c.JSON(200, base.ErrorDefault("UpLoad file error"))
		return
	}
	c.JSON(200, base.SuccessWithData(filePath))
}

// Download the file, when the file is not directly accessible, you need to use this interface to access
func (File) DownLoadFile(c *gin.Context) {
	// todo 从数据库获取文件地址信息

	err := leowebgin.DownLoadFile(c, "file")
	if err != nil {
		c.JSON(200, base.ErrorDefault(err.Error()))
		return
	}
}

func (File) UploadFiles(c *gin.Context) {
	filePath, err := leowebgin.UploadFiles(c)
	if err != nil {
		c.JSON(200, base.ErrorDefault("UpLoad file error"))
		return
	}
	c.JSON(200, base.SuccessWithData(filePath))
}

func (File) upload() func(c *gin.Context) {
	return func(c *gin.Context) {

	}
}
