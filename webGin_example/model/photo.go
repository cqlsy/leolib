package model

import (
	"github.com/cqlsy/leolib/leodb/mog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const colNamePhoto = "Photo"

type Photo struct {
	mog.BaseCol `json:"-" bson:"-"`
	Id          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	Cover       string               `json:"cover" bson:"cover"`
	Data        string               `json:"data" bson:"data"`
	CreateTime  int64                `json:"createTime" bson:"createTime"`
	UpdateTime  int64                `json:"updateTime" bson:"updateTime"`
	Order       int16                `json:"order" bson:"order"`
	Status      int8                 `json:"status" bson:"status"`   // 0： -1：delete
	OwnerId     primitive.ObjectID   `json:"ownerId" bson:"ownerId"` //
	AlbumId     []primitive.ObjectID `json:"albumId" bson:"albumId"` // 所属相册列表
	OriTime     int64               `json:"oriTime"  bson:"oriTime"`
}

func (Photo) CreateOrUpdate(id, user string, albums []string, name, data, cover, description string,
	order int16, status int8,oriTime int64) error {
	photoId := primitive.NilObjectID
	var err error
	if id != "" {
		photoId, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
	}
	userId, err := primitive.ObjectIDFromHex(user)
	if err != nil {
		return err
	}
	var albumIds []primitive.ObjectID
	for _, item := range albums {
		itemId, err := primitive.ObjectIDFromHex(item)
		if err != nil {
			continue
		}
		albumIds = append(albumIds, itemId)
	}
	if len(albumIds) == 0 {
		albumIds = append(albumIds, primitive.NilObjectID)
	}
	if order > 999 {
		order = 999
	}
	if order < 1 {
		order = 1
	}

	var photo = Photo{Id: photoId, OwnerId: userId, AlbumId: albumIds,
		Name: name, Data: data, Cover: cover, Order: order, Status: status,
		Description: description, OriTime: oriTime}
	if photo.Id == primitive.NilObjectID {
		//album.Id = primitive.NewObjectID()
		_, err := mongo.InsertOne(colNamePhoto, photo)
		return err
	}
	return mongo.UpdateOne(colNamePhoto, photo, "_id")
}

func (photo Photo) Update() error {
	photo.UpdateTime = time.Now().Unix()
	return mongo.UpdateOne(colNamePhoto, photo, "_id")
}

func (Photo) InsertMany(data []Photo) error {
	insert, err := mog.ToColArray(data)
	if err != nil {
		return err
	}
	_, err = mongo.InsertMany(colNamePhoto, insert)
	return err
}

func (Photo) FindByAlbum(userId, albumId string, page, pageSize int64) ([]Photo, int64, error) {
	filter := mog.NewFilter()
	result := make([]Photo, 0)
	id, err := primitive.ObjectIDFromHex(albumId)
	if err != nil {
		return result, 0, err
	}
	ownerId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return result, 0, err
	}
	filter.AddArrayMatch("albumId", id).
		AddEqualCondition("status", 0).
		AddEqualCondition("ownerId", ownerId)
	err = mongo.FindMany(&result, colNamePhoto, nil, map[string]int{"order": -1},
		filter, (page-1)*pageSize, pageSize)
	if err != nil {
		return result, 0, err
	}
	count, err := mongo.Count(colNamePhoto, filter)
	return result, count, err
}

func (Photo) FindOneById(id string) (Photo, error) {
	var result Photo
	photoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	err = mongo.FindOne(&result, colNamePhoto, nil, nil,
		mog.NewFilter().AddEqualCondition("_id", photoId))
	return result, err
}
