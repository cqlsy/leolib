package reqparams

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/cqlsy/leolib/leocrypto/aes"
	"github.com/cqlsy/leolib/leocrypto/rsa"
	"github.com/cqlsy/leolib/leowebgin"
	"github.com/cqlsy/leolib/webGin_example/server/base"
	"github.com/cqlsy/leolib/webGin_example/server/middle"
)

func GetReqPramsRsa(c *gin.Context) (map[string]interface{}, error) {
	data, err := leowebgin.GetReqParamsFromBody(c)
	if err != nil {
		c.JSON(200, base.ErrorDefault(err.Error()))
		return nil, err
	}
	str := ""
	switch data["data"].(type) {
	case string:
		str = data["data"].(string)
	default:
		//c.JSON(200, base.ErrorDefault("No Params"))
		return nil, nil
	}
	params := map[string]interface{}{}
	oriData, err := rsa.RsaDecrypt(middle.RsaPath+"private.pem", []byte(str))
	if err != nil {
		c.JSON(200, base.ErrorDefault("Params error"))
		return nil, err
	}
	err = json.Unmarshal(oriData, params)
	if err != nil {
		c.JSON(200, base.ErrorDefault("Params error for Get"))
		return nil, err
	}
	return params, nil
}

func GetReqPramsAes(c *gin.Context) (map[string]interface{}, error) {
	data, err := leowebgin.GetReqParamsFromBody(c)
	if err != nil {
		c.JSON(200, base.ErrorDefault(err.Error()))
		return nil, err
	}
	str := ""
	switch data["data"].(type) {
	case string:
		str = data["data"].(string)
	default:
		//c.JSON(200, base.ErrorDefault("No Params"))
		return nil, nil
	}
	key := ""
	d, exit := c.Get("aeskey")
	if !exit {
		c.JSON(200, base.ErrorDefault("Get User Info Err"))
		return nil, errors.New("Params error for Get ")
	}
	switch d.(type) {
	case string:
		key = data["aeskey"].(string)
	}
	if key == "" {
		c.JSON(200, base.ErrorDefault("Params error for Get"))
		return nil, errors.New("Params error for Get ")
	}
	params := map[string]interface{}{}
	oriData, err := aes.AesDecrypt([]byte(key), []byte(str))
	if err != nil {
		c.JSON(200, base.ErrorDefault("Params error"))
		return nil, err
	}
	err = json.Unmarshal(oriData, params)
	if err != nil {
		c.JSON(200, base.ErrorDefault("Params error for Get"))
		return nil, err
	}
	return params, nil
}

func GetReqPrams(c *gin.Context) (map[string]interface{}, error) {
	data, err := leowebgin.GetReqParamsFromBody(c)
	if err != nil {
		c.JSON(200, base.ErrorDefault(err.Error()))
		return nil, err
	}
	str := ""
	switch data["data"].(type) {
	case string:
		str = data["data"].(string)
	default:
		c.JSON(200, base.ErrorDefault("No Params"))
		return nil, err
	}
	params := map[string]interface{}{}
	err = json.Unmarshal([]byte(str), params)
	if err != nil {
		c.JSON(200, base.ErrorDefault("Params error for Get"))
		return nil, err
	}
	return params, nil
}
