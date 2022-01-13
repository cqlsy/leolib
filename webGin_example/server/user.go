package server

import (
	"github.com/gin-gonic/gin"
	"github.com/cqlsy/leolib/webGin_example/server/base"
	"github.com/cqlsy/leolib/webGin_example/server/middle"
	"github.com/cqlsy/leolib/webGin_example/server/reqparams"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
}

const DefaultUserId = "000000000000000000000001"

func GetDefaultUserId() primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(DefaultUserId)
	return id
}

func (User) Test() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, base.SuccessWithData("this"))
		//model.User{}.UpdateOne()
	}
}

func (User) GetRsaPubKey() func(c *gin.Context) {
	return func(c *gin.Context) {
		//rsa.GetPubKeyString(middle.PublicKey)
		if middle.PublicKey == "" {
			c.JSON(200, base.ErrorDefault("No Rsa PubKey"))
		} else {
			c.JSON(200, base.SuccessWithData(middle.PublicKey))
		}
	}
}

//
func (User) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		params, err := reqparams.GetReqPramsRsa(c)
		if err != nil || len(params) == 0 {
			return
		}
		println(params)
		// params[]

	}
}
