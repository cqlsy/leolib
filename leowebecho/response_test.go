/**
 * Created by WillkYang on 2017/3/10.
 */

package leowebecho

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
)

func mockResponse(url string) *Response {
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, url, strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	re := NewRequest(c)
	return NewResponse(c, re)
}
//
//func TestResponse_Success(t *testing.T) {
//	resp := mockResponse("/")
//	Equal(resp.Success("成功"), nil)
//}
//
//func TestResponse_SetStatus(t *testing.T) {
//	resp := mockResponse("/")
//	resp.SetStatus(500)
//	Equal(resp.Success("成功"), nil)
//	Equal(resp.Context().Response().Status, 500)
//	resp = mockResponse("/")
//	Equal(resp.Success("成功"), nil)
//	Equal(resp.Context().Response().Status, 200)
//}
