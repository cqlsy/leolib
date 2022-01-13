package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/cqlsy/leolib/leostrconv"
	"github.com/cqlsy/leolib/leoutil"
	"github.com/cqlsy/leolib/leowebgin"
	"github.com/cqlsy/leolib/webGin_example/model"
	"github.com/cqlsy/leolib/webGin_example/server/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type Album struct {
}

func (album Album) Init(engine *gin.Engine) {
	_ = model.InitDefault()
	group := engine.Group("/album")
	group.POST("/create", album.Create())
	group.GET("/maxOrder", album.GetMaxOrder())
	group.POST("/list", album.GetAlbums())
	group.POST("/update", album.Update())
	group.POST("/delete", album.delete())
}

func (Album) delete() gin.HandlerFunc {
	return func(context *gin.Context) {
		params, err := leowebgin.GetReqParamsFromBody(context)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		idStr := leoutil.GetString(params["id"])
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			context.JSON(200, base.Error(-1, "相册获取错误"))
			return
		}
		var album = model.Album{Id: id, Status: 1}
		err = album.UpdateAlbum()
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		context.JSON(200, base.SuccessDefault())
	}
}

func (Album) Update() gin.HandlerFunc {
	return func(context *gin.Context) {
		filePath, err := leowebgin.UploadFile(context, "cover")
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		idStr := leowebgin.GetReqParams(context, "id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		var album = model.Album{Id: id, Cover: filePath}
		name := leowebgin.GetReqParams(context, "name")
		description := leowebgin.GetReqParams(context, "description")
		order := leostrconv.AtoIDefault(leowebgin.GetReqParams(context, "order"), -1)
		status := leostrconv.AtoIDefault(leowebgin.GetReqParams(context, "status"), -1)
		//leowebgin.GetReqParams(context)
		if name != "" {
			album.Name = name
		}
		if description != "" {
			album.Description = description
		}
		if order != -1 {
			album.Order = int16(order)
		}
		if status != -1 {
			album.Status = int8(status)
		}
		err = album.UpdateAlbum()
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		context.JSON(200, base.SuccessDefault())
	}
}

func (Album) Create() gin.HandlerFunc {
	return func(context *gin.Context) {
		filePath, err := leowebgin.UploadFile(context, "cover")
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		name := leowebgin.GetReqParams(context, "name")
		description := leowebgin.GetReqParams(context, "description")
		order := leostrconv.AtoIDefault0(leowebgin.GetReqParams(context, "order"))
		status, _ := strconv.ParseInt(leowebgin.GetReqParams(context, "status"), 10, 8)
		ownerId := GetDefaultUserId()
		//leowebgin.GetReqParams(context)
		var album model.Album = model.Album{Name: name, Cover: filePath, Description: description,
			Order: int16(order), Status: int8(status), OwnerId: ownerId}
		err = album.CreateOrUpdateAlbum()
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		context.JSON(200, base.SuccessDefault())
	}
}

func (Album) GetMaxOrder() gin.HandlerFunc {
	return func(context *gin.Context) {
		var album model.Album
		order, err := album.FindMaxOrder(DefaultUserId)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
		} else {
			context.JSON(200, base.SuccessWithData(map[string]interface{}{"order": order}))
		}
	}
}

func (Album) GetAlbums() gin.HandlerFunc {
	return func(context *gin.Context) {
		params, err := leowebgin.GetReqParamsFromBody(context)
		if err != nil {
			context.JSON(200, base.ErrorDefault(err.Error()))
			return
		}
		data, count := model.Album{}.FindMany(DefaultUserId, leoutil.GetString(params["keyword"]), leoutil.GetInt64WithDefault(params["page"], 1),
			leoutil.GetInt64WithDefault(params["pageSize"], 10))
		for index, item := range data {
			if item.Cover != "" {
				tmp := item
				s := fmt.Sprintf("%v%v", base.GetPicDoMain(), item.Cover)
				tmp.Cover = s
				data[index] = tmp
			}
		}
		var result = map[string]interface{}{"data": data, "count": count}
		context.JSON(200, base.SuccessWithData(result))
	}
}
