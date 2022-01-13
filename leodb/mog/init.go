package mog

import (
	"context"
	"errors"
	"fmt"
	"github.com/cqlsy/leolib/leodb/config"
	"github.com/cqlsy/leolib/leolog"
	"github.com/cqlsy/leolib/leoutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"reflect"
	"time"
)

type MongoDb struct {
	db   *mongo.Database
	conf *config.Info
}

type Col interface {
	GetFilterKey(data interface{}, key string) interface{}
}

type test struct {
	BaseCol     `json:"-" bson:"-"`
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Cover       string             `json:"cover" bson:"cover"`
	CreateTime  int64              `json:"createTime" bson:"createTime"`
	UpdateTime  int64              `json:"updateTime" bson:"updateTime"`
	Order       int16              `json:"order" bson:"order"`
	Status      int8               `json:"status" bson:"status"` // 0：正常 -1：删除
}

type BaseCol struct {
	//data map[string]interface{} `json:"-" bson:"-"`
	//err  error                  `json:"-" bson:"-"`
}

func (base BaseCol) GetFilterKey(data interface{}, key string) interface{} {
	/*if base.data != nil && len(base.data) > 0 && base.err == nil {
		return base.data[key]
	}
	if base.err != nil {
		return nil
	}
	base.data, base.err = leoutil.Struct2Map(data)
	if base.err != nil {
		return nil
	}
	if value, ok := base.data[key]; ok {
		return value
	}*/
	value, _ := leoutil.GetStructValueByJsonId(data, key)
	return value
}

func Connect(Conf *config.Info) (*MongoDb, error) {
	ctx, cancel := NewContext(2)
	defer cancel()
	mongoInfo := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		Conf.Db.User,
		url.QueryEscape(Conf.Db.Password),
		Conf.Db.Host,
		Conf.Db.Port,
		Conf.Db.Db)
	clientOpts := options.Client().ApplyURI(mongoInfo)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	db := new(MongoDb)
	db.db = client.Database(Conf.Db.Db)
	leolog.Print("mongo connect success: " + mongoInfo)
	return db, nil //
}

func NewContextDEfault() (context.Context, func()) {
	return NewContext(10)
}

func NewContext(timeout time.Duration) (context.Context, func()) {
	return context.WithTimeout(context.Background(), timeout*time.Second)
}

func ToColArray(src interface{}) ([]Col, error) {
	var result []Col
	switch reflect.TypeOf(src).Kind() {
	case reflect.Slice, reflect.Array:
		ori := reflect.ValueOf(src)
		for i := 0; i < ori.Len(); i++ {
			item := ori.Index(i)
			switch item.Interface().(type) {
			case Col:
				result = append(result, item.Interface().(Col))
			default:
				return nil, errors.New("type error")
			}
		}
	}
	return result, nil
}
