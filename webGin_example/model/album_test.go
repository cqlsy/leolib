package model

import (
	"fmt"
	"github.com/cqlsy/leolib/leodb"
	"github.com/cqlsy/leolib/leoutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func initDB() {
	database := leodb.InitDataBase("../../config.json")
	Init(database)
}

func TestInit(t *testing.T) {
	//var albums []Album
	//albums = append(albums, Album{Name: "1", Description: "sss", Cover: "dsada"})
	//albums = append(albums, Album{Name: "2", Description: "ssss", Cover: "dsada"})
	//albums = append(albums, Album{Name: "3", Description: "sssd", Cover: "dsada"})
	//albums = append(albums, Album{Name: "4", Description: "sssf", Cover: "dsada"})
	//_, err := mog.ToColArray(albums)
	//if err != nil {
	//	println(err.Error())
	//} else {
	//
	//	println(fmt.Sprintf("%s", leoutil.ObjectToJSONStr(albums)))
	//}

	initDB()
	Album{Name: "1", Description: "sss", Cover: "dsada"}.CreateOrUpdateAlbum()
}

func TestAlbum_FindMany(t *testing.T) {
	initDB()
	data, count := Album{}.FindMany("000000000000000000000001", "", 1, 10)
	println(fmt.Sprintf("%v===%d", leoutil.ObjectToJSONStr(data), count))
}

func TestInitDefault(t *testing.T) {
	initDB()
	InitDefault()
}

func TestAlbum_FindMaxOrder(t *testing.T) {
	initDB()
	println(primitive.NilObjectID.Hex())
	data, err := Album{}.FindMaxOrder(primitive.NilObjectID.Hex())
	println(fmt.Sprintf("%v:::%v", err, data))
}

func TestAlbum_UpdateAlbum(t *testing.T) {
	initDB()
	id, _ := primitive.ObjectIDFromHex("61a4a139fea93bde9b741bf6")
	album := Album{Status: 5, Id: id, Name: "dsadasdsa", OwnerId: id}
	err := album.UpdateAlbum()
	if err != nil {
		println(err.Error())
	}
}
