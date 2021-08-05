package server

import (
	"github.com/gin-gonic/gin"
)

type Test struct {
}

func (Test) NetTest(c *gin.Context) {
	params := make(map[string]interface{})
	err := c.ShouldBind(params)
	//c.JSON(200, ResponseMsg{Code: 20000, Msg: "ok", Data: ""})
	var d int
	for i := 3; i > 1; i-- {
		d = 2 / i
	}
	if err != nil {
		//log.Print("ssssssssss")
	}
	//c.Header()
	//c.File()
	c.JSON(200, ResponseMsg{Code: 20000, Msg: "ok", Data: d})
}
