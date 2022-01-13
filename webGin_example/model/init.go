package model

import (
	"github.com/cqlsy/leolib/leodb"
	"github.com/cqlsy/leolib/leodb/mog"
)

var db *leodb.Db
var mongo *mog.MongoDb

func Init(leoDb *leodb.Db) {
	db = leoDb
	mongo = db.MogDb
}
