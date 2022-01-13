package model

import (
	"fmt"
	"github.com/cqlsy/leolib/leoutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestCreateOrUpdate(t *testing.T) {
	initDB()
	id, _ := primitive.ObjectIDFromHex("61a442189168322c3bece8a3")
	err := Photo{}.CreateOrUpdate("", primitive.NilObjectID.Hex(),
		[]string{
			primitive.NilObjectID.Hex(),
			id.Hex(),
		}, "Default",
		"data", "cover", "description", 12, 0, time.Now().Unix())
	if err != nil {
		println(err.Error())
	}
}

func TestPhoto_FindByAlbum(t *testing.T) {
	initDB()
	result, count, err := Photo{}.FindByAlbum(primitive.NilObjectID.Hex(),
		"61a442189168322c3bece8a4", 1, 10)

	println(fmt.Sprintf("%v:::%v:::%v", err, leoutil.ObjectToJSONStr(result), count))
}

func TestPhoto_FindOneById(t *testing.T) {
	initDB()
	result, err := Photo{}.FindOneById("61a442189168322c3bece8a3")

	println(fmt.Sprintf("%v:::%v", err, leoutil.ObjectToJSONStr(result)))
}
