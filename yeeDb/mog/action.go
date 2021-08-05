package mog

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	config "github.com/lsy/yeelib/conf"
	"net/url"
	"time"
)

type MongoDb struct {
	Db  *mongo.Database
	ctx context.Context
}

// 链接数据库
func Collect() (*MongoDb, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	mongoInfo := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		config.Conf.Db.User,
		url.QueryEscape(config.Conf.Db.Password),
		config.Conf.Db.Host,
		config.Conf.Db.Port,
		config.Conf.Db.Db)
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
	db.ctx = context.Background()
	db.Db = client.Database(config.Conf.Db.Db)
	return db, nil // 链接到我们需要的数据库
}

/**
更新数据，没有的数据话插入数据
*/
func (db *MongoDb) UpdateOrInsertOne(colName string, document map[string]interface{}, filterKey string) error {
	if len(filterKey) == 0 {
		return errors.New("UpdateOrInsert Must has more than one filterKey")
	}
	filter := bson.D{}
	if len(filterKey) > 0 {
		filter = bson.D{{filterKey, document[filterKey]}}
	}
	update := bson.D{{"$set", document}}
	updateOpts := options.Update().SetUpsert(true) // 如果没有就插入新的数据
	_, err := db.Db.Collection(colName).UpdateOne(db.ctx, filter, update, updateOpts)
	return err
}

/**
更新数据，没有的数据话插入数据
*/
func (db *MongoDb) UpdateOrInsertOneWithKeys(colName string, document map[string]interface{}, filterKeys []string) error {
	var hasKey = false
	filter := bson.D{}
	for _, key := range filterKeys {
		if len(key) > 0 {
			filter = append(filter, bson.E{Key: key, Value: document[key]})
			hasKey = true
		}
	}
	if !hasKey {
		return errors.New("UpdateOrInsert Must has more than one filterKey")
	}
	update := bson.D{{"$set", document}}
	updateOpts := options.Update().SetUpsert(true) // 如果没有就插入新的数据
	_, err := db.Db.Collection(colName).UpdateOne(db.ctx, filter, update, updateOpts)
	return err
}

func (db *MongoDb) InsertMany(document []map[string]interface{}, colName string) error {
	return db.UpdateOrInsertMany(colName, document, "")
}

/**
更新数据，没有的数据话插入数据
*/
func (db *MongoDb) UpdateOrInsertMany(colName string, document []map[string]interface{}, filterKey string) error {
	if len(filterKey) == 0 {
		return errors.New("UpdateOrInsert Must has more than one filterKey")
	}
	var models = make([]mongo.WriteModel, 0)
	for _, value := range document {
		filter := bson.D{}
		if len(filterKey) > 0 {
			filter = bson.D{{filterKey, value[filterKey]}}
		}
		models = append(models, mongo.NewUpdateOneModel().
			SetUpsert(true).
			SetFilter(filter).
			SetUpdate(bson.D{{"$set", value}}))
	}
	opts := options.BulkWrite().SetOrdered(false)
	_, err := db.Db.Collection(colName).BulkWrite(db.ctx, models, opts)
	// 这里不执行任何操作，
	return err
}

/**
更新数据，没有的数据话插入数据
*/
func (db *MongoDb) UpdateOrInsertManyWithKeys(colName string, document []map[string]interface{}, filterKeys []string) error {
	var models = make([]mongo.WriteModel, 0)
	var hasKey = false
	for _, value := range document {
		filter := bson.D{}
		for _, key := range filterKeys {
			if len(key) > 0 {
				filter = append(filter, bson.E{Key: key, Value: value[key]})
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
	_, err := db.Db.Collection(colName).BulkWrite(db.ctx, models, opts)
	// 这里不执行任何操作，
	return err
}

/// 数据的查找，这里只是用一些例子记录，方便今后的查询
// sortMap {"name",1(-1)}  排序的字段 1：正序 -1 倒序 不需要传入 map[string]int
// projects 需要查询的字段 ，全部查询： []string{}
// filter 过滤字段 ，不需要传入nil
func (db *MongoDb) FindOne(colName string, projects []string, sortMap map[string]int, filter *Filter) (map[string]interface{}, error) {
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
	result := make(map[string]interface{})
	err := db.Db.Collection(colName).FindOne(db.ctx, f, findOptions).Decode(result)
	if err != nil {
		// 获取失败
		return nil, err
	}
	return result, err
}

// 查询数据总数
func (db *MongoDb) Count(colName string, filter *Filter) (int64, error) {
	opts := options.Count()
	f := bson.D{}
	if filter != nil {
		f = filter.value
	}
	txCount, err := db.Db.Collection(colName).CountDocuments(db.ctx, f, opts)
	if err != nil {
		// 获取失败
		return 0, err
	}
	return txCount, err
}

// sortMap {"name",1(-1)} 排序的字段 1：正序 -1 倒序 不需要传入 map[string]int
// projects 需要查询的字段 ，全部查询： []string{}
// filter 过滤字段 ，不需要传入nil
func (db *MongoDb) FindMany(colName string, projects []string, sortMap map[string]int, filter *Filter, skip int64, limit int64) ([]map[string]interface{}, error) {
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
	result := make([]map[string]interface{}, 0)
	c, err := db.Db.Collection(colName).Find(db.ctx, f, findOptions)
	if c == nil || err != nil {
		return nil, err
	}
	err = c.All(db.ctx, &result)
	if err != nil {
		// 获取失败
		return nil, err
	}
	return result, err
}

// 最终的Filter
type Filter struct {
	value bson.D
}

// Or条件语句, 只能添加Filter
type OrFilter struct {
	value []bson.D
}

func NewFilter() *Filter {
	result := new(Filter)
	result.value = bson.D{}
	return result
}

func NewOrFilter() *OrFilter {
	result := new(OrFilter)
	result.value = []bson.D{}
	return result
}

// 增加单条的or条件查询语句
func (or *OrFilter) AddCondition(filter *Filter) {
	or.value = append(or.value, filter.value)
}

// 把 or条件查询语句加入到过滤语句中
func (fil *Filter) AddOrCondition(or *OrFilter) {
	if len(or.value) <= 0 {
		return
	}
	fil.value = append(fil.value, bson.E{Key: "$or", Value: or.value})
}

func (fil *Filter) AddEqualCondition(key string, value interface{}) {
	fil.value = append(fil.value, bson.E{Key: key, Value: value})
}

func (fil *Filter) AddNeCondition(key string, value interface{}) {
	fil.value = append(fil.value, bson.E{Key: key, Value: bson.D{{Key: "$ne", Value: value}}})
}

//  primitive.Regex{Pattern: params.Keywords}} like , 模糊查询
func (fil *Filter) AddLikeCondition(key string, value string) {
	//i 表示不区分大小写
	fil.value = append(fil.value, bson.E{Key: key, Value: primitive.Regex{Pattern: value, Options: "i"}})
}

// 字段是否存在
func (fil *Filter) AddExistsCondition(key string, value bool) {
	fil.value = append(fil.value, bson.E{Key: key, Value: bson.D{{Key: "$exists", Value: value}}})
}

func (fil *Filter) AddBiggerCondition(key string, value interface{}, isEqual bool) {
	str := "$gt"
	if isEqual {
		str = "$gte"
	}
	fil.value = append(fil.value, bson.E{Key: key, Value: bson.D{{Key: str, Value: value}}})
}

func (fil *Filter) AddLessCondition(key string, value interface{}, isEqual bool) {
	str := "$lt"
	if isEqual {
		str = "$lte"
	}
	fil.value = append(fil.value, bson.E{Key: key, Value: bson.D{{Key: str, Value: value}}})
}
