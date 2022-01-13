package model

import "time"

type ColUser struct {
	/**
	name
	nickName
	phone
	email
	password
	lastLoginTime
	lastActiveTime
	token
	*/
}

const colNameUser = "user"
const TokenExpireTime = 60 * 60 * 24 * 7 // unit second

// user ==> get
func (ColUser) Register(phone string, name string, password string) error {
	data := make(map[string]interface{})
	data["phone"] = phone
	if name == "" {
		name = "Leo"
	}
	data["name"] = name
	data["password"] = password
	data["registerTime"] = time.Now().Second()
	//err := db.MogDb.UpdateOne(colName, data, "phone")
	//return err
	return nil
}
