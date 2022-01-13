package mog

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Filter struct {
	value bson.D
}

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

func (or *OrFilter) AddFilter(filter *Filter) *OrFilter {
	or.value = append(or.value, filter.value)
	return or
}

func (fil *Filter) AddOrOrFilter(or *OrFilter) *Filter {
	if len(or.value) <= 0 {
		return fil
	}
	//fil.value = append(fil.value, bson.E{Key: "$or", Value: or.value})
	fil.AddBsonE(bson.E{Key: "$or", Value: or.value})
	return fil
}

func (fil *Filter) AddEqualCondition(key string, value interface{}) *Filter {
	//fil.value = append(fil.value, bson.E{Key: key, Value: value})
	fil.AddBsonE(bson.E{Key: key, Value: value})
	return fil
}

func (fil *Filter) AddNeCondition(key string, value interface{}) *Filter {
	//fil.value = append(fil.value, bson.E{Key: key, Value: bson.D{{Key: "$ne", Value: value}}})
	fil.AddBsonE(bson.E{Key: key, Value: bson.D{{Key: "$ne", Value: value}}})
	return fil
}

func (fil *Filter) AddLikeCondition(key string, value string) *Filter {
	//fil.value = append(fil.value, bson.E{Key: key, Value: primitive.Regex{Pattern: value, Options: "i"}})
	fil.AddBsonE(bson.E{Key: key, Value: primitive.Regex{Pattern: value, Options: "i"}})
	return fil
}

func (fil *Filter) AddExistsCondition(key string, value bool) *Filter {
	//fil.value = append(fil.value, bson.E{Key: key, Value: bson.D{{Key: "$exists", Value: value}}})
	fil.AddBsonE(bson.E{Key: key, Value: bson.D{{Key: "$exists", Value: value}}})
	return fil
}

func (fil *Filter) AddBiggerCondition(key string, value interface{}, isEqual bool) *Filter {
	str := "$gt"
	if isEqual {
		str = "$gte"
	}
	fil.AddBsonE(bson.E{Key: key, Value: bson.D{{Key: str, Value: value}}})
	return fil
}

func (fil *Filter) AddLessCondition(key string, value interface{}, isEqual bool) *Filter {
	str := "$lt"
	if isEqual {
		str = "$lte"
	}
	fil.AddBsonE(bson.E{Key: key, Value: bson.D{{Key: str, Value: value}}})
	return fil
}

func (fil *Filter) AddBsonE(condition bson.E) *Filter {
	fil.value = append(fil.value, condition)
	return fil
}

// group
/**
bson.D{
		{"$group", bson.D{
			{"_id", null},
			{"balance", bson.D{
				{"$sum", "$balance"},
			}},
		}},
*/
// https://blog.csdn.net/qq_18948359/article/details/88777066 for more info
// key group by this column
// conditions can no set
func (fil *Filter) AddGroupFilter(key string, conditions []bson.E) *Filter {
	id := "null"
	if len(key) > 0 {
		id = "$" + key
	}
	filter := bson.D{
		{"_id", id},
	}
	for _, item := range conditions {
		filter = append(filter, item)
	}
	filValue := bson.E{Key: "$group", Value: filter}
	fil.AddBsonE(filValue)
	return fil
}

// 数组中的筛选查询 "$elemMatch"
func (fil *Filter) AddArrayMatch(key string, value interface{}) *Filter {
	fil.AddBsonE(
		bson.E{Key: key, Value:
		bson.M{"$elemMatch": bson.M{"$gte": value, "$lte": value}}},
	)
	return fil
}
