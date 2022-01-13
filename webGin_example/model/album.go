package model

import (
	"errors"
	"github.com/cqlsy/leolib/leodb/mog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"time"
)

const colNameAlbum = "Album"

type Album struct {
	mog.BaseCol `json:"-" bson:"-"`
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	Cover       string             `json:"cover" bson:"cover,omitempty"`
	CreateTime  int64              `json:"createTime" bson:"createTime,omitempty"`
	UpdateTime  int64              `json:"updateTime" bson:"updateTime,omitempty"`
	Order       int16              `json:"order" bson:"order,omitempty"`
	Status      int8               `json:"status" bson:"status,omitempty"` // 0：1：delete
	OwnerId     primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
}

func InitDefault() error {
	var result []Album
	err := mongo.FindMany(&result, colNameAlbum, []string{"ownerId"}, nil, mog.NewFilter().
		AddEqualCondition("ownerId", primitive.NilObjectID), 0, math.MaxInt64)
	if err != nil || len(result) == 0 {
		var del = Album{Name: "默认", Description: "这是第一个默认相册，当你上传的图片未指定的时候，系统会将其归纳到默认相册",
			Order: 1000, CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}
		var albums []Album
		albums = append(albums, del)
		insert, err := mog.ToColArray(albums)
		if err != nil {
			return err
		}
		_, err = mongo.InsertMany(colNameAlbum, insert)
		return err
	}
	return nil
}

func (album Album) CreateOrUpdateAlbum() error {
	album.UpdateTime = time.Now().Unix()
	if album.Id == primitive.NilObjectID {
		//album.Id = primitive.NewObjectID()
		album.CreateTime = time.Now().Unix()
		_, err := mongo.InsertOne(colNameAlbum, album)
		return err
	} else {
		err := mongo.UpdateOne(colNameAlbum, album, "_id")
		return err
	}
}

func (album Album) UpdateAlbum() error {
	album.UpdateTime = time.Now().Unix()
	err := mongo.UpdateOne(colNameAlbum, album, "_id")
	return err
}

func (album Album) FindMany(ownerId, keyword string, page, pageNum int64) ([]Album, int64) {
	data := make([]Album, 0)
	id, err := primitive.ObjectIDFromHex(ownerId)
	if err != nil {
		return data, 0
	}
	filter := mog.NewFilter().AddEqualCondition("status", 0).AddOrOrFilter(
		mog.NewOrFilter().
			AddFilter(mog.NewFilter().AddEqualCondition("ownerId", id)).
			AddFilter(mog.NewFilter().AddEqualCondition("ownerId", primitive.NilObjectID))).
		AddOrOrFilter(mog.NewOrFilter().
			AddFilter(mog.NewFilter().AddLikeCondition("name", keyword)).
			AddFilter(mog.NewFilter().AddLikeCondition("description", keyword)))
	err = mongo.FindMany(&data, colNameAlbum, nil,
		map[string]int{"order": -1, "createTime": -1},
		filter, pageNum*(page-1), pageNum)
	if err != nil {
		data = make([]Album, 0)
	}
	count, err := mongo.Count(colNameAlbum, filter)
	if err != nil {
		count = int64(len(data))
	}
	return data, count
}

func (album Album) FindMaxOrder(ownerId string) (int16, error) {
	id, err := primitive.ObjectIDFromHex(ownerId)
	if err != nil {
		return 0, err
	}
	if id == primitive.NilObjectID {
		return 0, errors.New("invalid Id")
	}
	var resilt Album
	err = mongo.FindOne(&resilt, colNameAlbum, []string{"order"}, map[string]int{"order": -1},
		mog.NewFilter().AddEqualCondition("ownerId", id))
	if err != nil {
		return 0, err
	}

	return resilt.Order, nil
}
