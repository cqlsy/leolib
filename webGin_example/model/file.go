package model

type File struct {
}

var colNameFile = "File"

/**
type:article,image/png,video（文件类型）
path:serve path（服务器文件地址，用于访问）
name:name（上传的文件名称）
recordTime:（记录文件的时间）
uploadTime:（更新时间）
author:（上传者）
memo:（备注信息）
memoName:（备注名称）
visit:（浏览量）
size:（文件大小）
*/

//
//func (File) UploadFiles(data []map[string]interface{}) error {
//	_, err := db.MogDb.InsertMany(colNameFile, data )
//	return err
//}
