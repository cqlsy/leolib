package mog

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *MongoDb) InsertOne(colName string, document Col) (string, error) {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	res, err := db.db.Collection(colName).InsertOne(ctx, document, options.InsertOne())
	result := ""
	if err == nil {
		result = res.InsertedID.(primitive.ObjectID).String()
	}
	return result, err
}

/**
 */
func (db *MongoDb) UpdateOne(colName string, document Col, filterKey string) error {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	if len(filterKey) == 0 {
		return errors.New("UpdateOne Must has more than one filterKey")
	}
	filter := bson.D{}
	if len(filterKey) > 0 {
		str := document.GetFilterKey(document, filterKey)
		filter = bson.D{{filterKey, str}}
	}
	update := bson.D{{"$set", document}}
	updateOpts := options.Update().SetUpsert(true)
	_, err := db.db.Collection(colName).UpdateOne(ctx, filter, update, updateOpts)
	return err
}

func (db *MongoDb) UpdateOrInsertOneWithKeys(colName string, document Col, filterKeys []string) error {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	var hasKey = false
	filter := bson.D{}
	for index, key := range filterKeys {
		if len(key) > 0 {
			filter = append(filter, bson.E{Key: key, Value: document.GetFilterKey(document, filterKeys[index])})
			hasKey = true
		}
	}
	if !hasKey {
		return errors.New("UpdateOrInsert Must has more than one filterKey")
	}
	update := bson.D{{"$set", document}}
	updateOpts := options.Update().SetUpsert(true)
	_, err := db.db.Collection(colName).UpdateOne(ctx, filter, update, updateOpts)
	return err
}

// document 是实现Col的列表数据// 或者map
func (db *MongoDb) InsertMany(colName string, document []Col) ([]string, error) {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	if document == nil || len(document) == 0 {
		return nil, errors.New("no data to save")
	}
	var doc = make([]interface{}, 0)
	for _, item := range document {
		doc = append(doc, item)
	}
	result, err := db.db.Collection(colName).InsertMany(ctx, doc, options.InsertMany())
	res := []string{}
	if err == nil && len(result.InsertedIDs) > 0 {
		for _, value := range result.InsertedIDs {
			switch value.(type) {
			case primitive.ObjectID:
				res = append(res, value.(primitive.ObjectID).String())
			default:
				res = append(res, fmt.Sprintf("%v", value))
			}
		}
	}
	return res, err
}

/**
 */
func (db *MongoDb) UpdateOrInsertMany(colName string, filterKey string,  document []Col) error {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	if len(filterKey) == 0 {
		return errors.New("UpdateOrInsert Must has more than one filterKey")
	}
	var models = make([]mongo.WriteModel, 0)
	for _, value := range document {
		filter := bson.D{}
		if len(filterKey) > 0 {
			filter = bson.D{{filterKey, value.(Col).GetFilterKey(document, filterKey)}}
		}
		models = append(models, mongo.NewUpdateOneModel().
			SetUpsert(true).
			SetFilter(filter).
			SetUpdate(bson.D{{"$set", value}}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := db.db.Collection(colName).BulkWrite(ctx, models, opts)
	return err
}

/**
 */
func (db *MongoDb) UpdateOrInsertManyWithKeys(colName string, filterKeys []string, document []Col) error {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	var models = make([]mongo.WriteModel, 0)
	var hasKey = false
	for _, value := range document {
		filter := bson.D{}
		for index, key := range filterKeys {
			if len(key) > 0 {
				filter = bson.D{{key, value.(Col).GetFilterKey(value, filterKeys[index])}}
				hasKey = true
			}
		}
		models = append(models, mongo.NewUpdateOneModel().
			SetUpsert(true).
			SetFilter(filter).
			SetUpdate(bson.D{{"$set", value}}))
	}
	if !hasKey {
		return errors.New("UpdateOrInsert Must has more than one filterKey")
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := db.db.Collection(colName).BulkWrite(ctx, models, opts)
	return err
}

func (db *MongoDb) FindOne(result Col, colName string, projects []string, sortMap map[string]int, filter *Filter) error {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	sort := bson.D{}
	for key, value := range sortMap {
		sort = append(sort, bson.E{Key: key, Value: value})
	}
	findOptions := options.FindOne()
	findOptions.SetSort(sort)
	if len(projects) > 0 {
		project := bson.D{}
		for _, value := range projects {
			project = append(project, bson.E{Key: value, Value: 1})
		}
		findOptions.SetProjection(project)
	}
	f := bson.D{}
	if filter != nil {
		f = filter.value
	}
	err := db.db.Collection(colName).FindOne(ctx, f, findOptions).Decode(result)
	if err != nil {
		return err
	}
	return err
}

func (db *MongoDb) Count(colName string, filter *Filter) (int64, error) {
	ctx, cancel := NewContextDEfault()
	defer cancel()
	opts := options.Count()
	f := bson.D{}
	if filter != nil {
		f = filter.value
	}
	txCount, err := db.db.Collection(colName).CountDocuments(ctx, f, opts)
	if err != nil {
		return 0, err
	}
	return txCount, err
}

func (db *MongoDb) FindMany(
	result interface{},
	colName string,
	projects []string,
	sortMap map[string]int,
	filter *Filter,
	skip int64, limit int64,
) error {
	//println(skip, limit)
	ctx, cancel := NewContextDEfault()
	defer cancel()
	sort := bson.D{}
	if sortMap != nil && len(sortMap) > 0 {
		for key, value := range sortMap {
			sort = append(sort, bson.E{Key: key, Value: value})
		}
	}
	findOptions := options.Find()
	findOptions.SetSort(sort)
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	if projects != nil && len(projects) > 0 {
		project := bson.D{}
		for _, value := range projects {
			project = append(project, bson.E{Key: value, Value: 1})
		}
		findOptions.SetProjection(project)
	}
	f := bson.D{}
	if filter != nil {
		f = filter.value
	}
	c, err := db.db.Collection(colName).Find(ctx, f, findOptions)
	if c == nil || err != nil {
		return err
	}
	err = c.All(ctx, result)
	if err != nil {
		return err
	}
	return err
}
