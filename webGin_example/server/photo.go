package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/cqlsy/leolib/leoutil"
	"github.com/cqlsy/leolib/leowebgin"
	"github.com/cqlsy/leolib/webGin_example/model"
	"github.com/cqlsy/leolib/webGin_example/server/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type Photo struct {
}

func (photo Photo) Init(engine *gin.Engine) {
	group := engine.Group("/photo")
	group.POST("/create", photo.create())
	group.POST("/update", photo.updateOne())
	group.POST("/list", photo.photoList())
	group.POST("/delete", photo.delete())
}

func (Photo) create() gin.HandlerFunc {
	return func(context *gin.Context) {
		filePaths, err := leowebgin.UploadFiles(context)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		var idStr []string
		albunStr := leowebgin.GetReqParams(context, "albumIds")
		if albunStr != "" {
			idStr = strings.Split(albunStr, ",")
		}
		var ids []primitive.ObjectID
		if len(idStr) == 0 {
			ids = append(ids, primitive.NilObjectID)
		} else {
			for _, item := range idStr {
				id, err := primitive.ObjectIDFromHex(item)
				if err != nil {
					context.JSON(200, base.ErrorDefault(err.Error()))
					return
				}
				ids = append(ids, id)
			}
		}
		oriTime := strings.Split(leowebgin.GetReqParams(context, "oriTime"), ",")
		var photos []model.Photo
		for index, item := range filePaths {
			photos = append(photos, model.Photo{Cover: item, Data: item, CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix(),
				Order: 1, OwnerId: GetDefaultUserId(), AlbumId: ids,
				OriTime: leoutil.GetInt64WithDefault(oriTime[index], time.Now().Unix())})
		}
		err = model.Photo{}.InsertMany(photos)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		context.JSON(200, base.SuccessDefault())
	}
}

func (Photo) delete() gin.HandlerFunc {
	return func(context *gin.Context) {
		params, err := leowebgin.GetReqParamsFromBody(context)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		idStr := leoutil.GetString(params["id"])
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			context.JSON(200, base.Error(-1, "图片获取错误"))
			return
		}
		var photo = model.Photo{Id: id, Status: 1}
		err = photo.Update()
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		context.JSON(200, base.SuccessDefault())
	}
}

// 更新排序，描述，名称
func (Photo) updateOne() gin.HandlerFunc {
	return func(context *gin.Context) {
		params, err := leowebgin.GetReqParamsFromBody(context)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		idStr := leoutil.GetString(params["id"])
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			context.JSON(200, base.Error(-1, "图片ID获取错误"))
			return
		}
		photo := model.Photo{Id: id}
		if value, ok := params["order"]; ok {
			photo.Order = int16(leoutil.GetInt64WithDefault(value, 1))
		}
		if value, ok := params["description"]; ok {
			photo.Description = leoutil.GetString(value)
		}
		if value, ok := params["name"]; ok {
			photo.Name = leoutil.GetString(value)
		}
		if value, ok := params["order"]; ok {
			photo.Status = int8(leoutil.GetInt64WithDefault(value, 0))
		}
		err = photo.Update()
		if err != nil {
			context.JSON(200, base.Error(-1, "更新错误"))
			return
		}
		context.JSON(200, base.SuccessDefault())
	}
}

//
func (Photo) photoList() gin.HandlerFunc {
	return func(context *gin.Context) {
		params, err := leowebgin.GetReqParamsFromBody(context)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		albumId := leoutil.GetString(params["albumId"])
		if albumId == "" {
			context.JSON(200, base.ErrorDefault("请选择相册"))
			return
		}
		data, count, err := model.Photo{}.FindByAlbum(DefaultUserId, albumId,
			leoutil.GetInt64WithDefault(params["page"], 1),
			leoutil.GetInt64WithDefault(params["pageSize"], 10))
		for index, item := range data {
			if item.Cover != "" {
				tmp := item
				s := fmt.Sprintf("%v%v", base.GetPicDoMain(), item.Cover)
				tmp.Cover = s
				tmp.Data = s
				data[index] = tmp
			}
		}
		var result = map[string]interface{}{"data": data, "count": count}
		context.JSON(200, base.SuccessWithData(result))
	}
}
