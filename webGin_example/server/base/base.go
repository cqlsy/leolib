package base

var picDoMain = "http://localhost:9999"

func GetPicDoMain() string {
	return picDoMain
}

func InitPicDoMain(domain string) {
	picDoMain = domain
}

type responseMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func Response(code int, msg string, data interface{}) responseMsg {
	return responseMsg{
		code,
		msg,
		data,
	}
}

func SuccessWithData(data interface{}) responseMsg {
	return Response(0, "Success", data)
}

func SuccessDefault() responseMsg {
	return Response(0, "Success", nil)
}

func SuccessMessage(msg string) responseMsg {
	return Response(0, msg, nil)
}

func LoginInvalid() responseMsg {
	return Response(-99, "Login invalid", nil)
}

func ErrorDefault(msg string) responseMsg {
	return Response(-1, msg, nil)
}

func Error(code int, msg string) responseMsg {
	if code == 0 {
		code = -1
	}
	return Response(code, msg, nil)
}
