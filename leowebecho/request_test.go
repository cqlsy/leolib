package leowebecho

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func mockRequest_Get(url string) *Request {
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, url, strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return NewRequest(c)
}

func mockRequest_Post(url string, values url.Values) *Request {
	e := echo.New()
	req, _ := http.NewRequest(echo.POST, url, strings.NewReader(""))
	req.Form = values
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return NewRequest(c)
}

//func TestNewRequest(t *testing.T) {
//	e := echo.New()
//	req, _ := http.NewRequest(echo.POST, "/", strings.NewReader(""))
//	rec := httptest.NewRecorder()
//	NotEqual(NewRequest(e.NewContext(req, rec)), nil)
//}
//
//func TestRequest_GetParam(t *testing.T) {
//	req := mockRequest_Get("/?key1=str&key2=2&key3=1.34534")
//	Equal(req.GetParam("no").SetDefault("yes").GetString(), "yes")
//	Equal(req.GetParam("111").GetString(), "")
//	Equal(req.GetParam("key1").GetString(), "str")
//	Equal(req.GetParam("key2").GetInt(), 2)
//	Equal(req.GetParam("key3").GetFloat(), 1.34534)
//}
//
//func TestRequest_PostParam(t *testing.T) {
//	req := mockRequest_Post("/", url.Values{
//		"key1": []string{"12312"},
//		"key2": []string{"3423.234"},
//	})
//	Equal(req.PostParam("key1").GetString(), "12312")
//	Equal(req.PostParam("key1").GetInt(), 12312)
//	Equal(req.PostParam("key2").GetFloat(), 3423.234)
//}
//
//func TestRequest_Validation(t *testing.T) {
//	req := mockRequest_Get("/?name=thisisaname&age=21" +
//		"&mobile=13111111111&tel=87221828&email=test@leoyuntech.com&postcode=672367&bash=_jlk2&ip=10.1.1.1")
//	Equal(req.Param("name").MinLength(1).GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").MinLength(12).GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").MaxLength(1).GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("name").MaxLength(12).GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("age").Min(1).GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("age").Min(100).GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("age").Max(1).GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("age").Max(100).GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("tel").Phone().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("mobile").Phone().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("mobile").Mobile().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("mobile").Tel().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("tel").Mobile().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("tel").Tel().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("email").Email().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").Email().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("tel").Email().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("postcode").ZipCode().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").ZipCode().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("tel").ZipCode().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("postcode").Numeric().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").Numeric().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("postcode").Alpha().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("name").Alpha().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("postcode").AlphaNumeric().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("bash").AlphaNumeric().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("bash").AlphaDash().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("ip").IP().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").IP().GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("bash").IP().GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("ip").Match("10").GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("name").Match("10").GetError(), nil)
//	req.CleanError()
//	NotEqual(req.Param("ip").NoMatch("10").GetError(), nil)
//	req.CleanError()
//	Equal(req.Param("name").NoMatch("10").GetError(), nil)
//	req.CleanError()
//}
//
//func TestRequest_SetJson(t *testing.T) {
//	req := mockRequest_Get("/")
//	req.SetJson(`["value",1]`)
//	Equal(req.JsonParam().GetJson().GetIndex(1).ToString(), "value")
//	Equal(req.JsonParam().GetJson().GetIndex(2).ToInt(), 1)
//}
//
//func TestRequest_GetJson(t *testing.T) {
//	req := mockRequest_Get("/")
//	req.SetJson(`{"value":"test"}`)
//	Equal(req.JsonParam("value").GetJson().ToString(), "test")
//}
//
//func TestRequest_SetDefault(t *testing.T) {
//	req := mockRequest_Get("/?value2=test22")
//	Equal(req.JsonParam("value1").SetDefault("test11").GetJson().ToString(), "test11")
//}
