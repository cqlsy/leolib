package leoconfig

var Conf Info

type Info struct {
	Db  db
	Web web
}

type web struct {
	LogPath      string
	RunMode      string
	Ip           string
	Port         int
	SaveFilePath string
	PicDoMain    string
}

type db struct {
	Host     string
	Db       string
	Port     string
	User     string
	Password string
	Protocol string // mongodb / mysql
}
