package leoutil

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

type User struct {
	Id   primitive.ObjectID
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
	Des  string `json:"des"`
}

func TestMap2Struct(t *testing.T) {
	d := map[string]interface{}{"name": "name", "age": 1, "des": "des"}
	var user User = User{}
	err := Map2Struct(d, &user)
	if err != nil {
		println("error", err.Error())
	} else {
		println(fmt.Sprintf("%v", d))
		println(fmt.Sprintf("%v", user))
	}

}

func TestStruct2Map(t *testing.T) {
	user := User{Id: primitive.NewObjectID(), Name: "name", Age: 2, Des: "des"}
	d, _ := Struct2Map(user)
	println(fmt.Sprintf("%v", d))
	oId, _ := primitive.ObjectIDFromHex("61a06b4e176681aa739814db")
	println(oId.String())
}

func TestStruct2MapD(t *testing.T) {
	users := User{Id: primitive.NewObjectID(), Name: "name", Age: 2, Des: "des"}
	resultsVal := reflect.ValueOf(users)
	var user interface{}
	if resultsVal.Kind() != reflect.Ptr {
		user = &users
	} else {
		user = users
	}
	types := reflect.TypeOf(user)
	//value := reflect.ValueOf(&user)
	name := ""
	for index := 0; index < types.Elem().NumField(); index++ {
		ss := types.Elem().Field(index).Tag.Get("json")
		println(ss)
		if ss == "age" {
			name = types.Elem().Field(index).Name
			break
		}
	}
	println(name)
	//data := value.FieldByName(name).Interface()
	//println(fmt.Sprintf("%v", data))
}

func TestGetStructValueByJsonId(t *testing.T) {
	users := User{Id: primitive.NewObjectID(), Name: "name", Age: 2, Des: "des"}
	s, _ := GetStructValueByJsonId(users, "age")
	println(fmt.Sprintf("%v",s))
}
