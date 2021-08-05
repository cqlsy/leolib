package server

import (
	"github.com/gin-gonic/gin"
	"lib/yeefile"
)

type File struct {
}

func (File) UploadFile(c *gin.Context) {
	//获取表单数据 参数为name值
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, ResponseMsg{Code: 20000, Msg: "Get file error"})
		return
	}
	filPath, err := yeefile.SaveFile(fileHeader)
	c.JSON(200, ResponseMsg{20000, "success", filPath})
}

// 下载文件，当文件不是直接能访问到，就需要使用这个接口来访问
func (File) DownLoadFile(c *gin.Context) {
	//filePath := c.Params.ByName("filePath")
	//filePath = c.Param("filePath")
	filePath := c.Query("filePath")   // 获取 Query（拼接到Url后面的参数）
	filePath = c.PostForm("filePath") // Post 的请求参数
	data := make(map[string]interface{})
	c.ShouldBind(&data)
	if !yeefile.FileExists(filePath) {
		c.JSON(200, ResponseMsg{-1, "error, no such file: " + filePath, ""})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+yeefile.GetFilName(filePath))
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}
