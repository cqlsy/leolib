package config

type Info struct {
	Db db
}

type db struct {
	Host        string
	Db          string
	Port        string
	User        string
	Password    string
	Protocol    string // mongodb / mysql
	Collections map[string]string
}
