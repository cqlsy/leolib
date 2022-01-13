/**
 * Created by angelina on 2017/4/15.
 */

package leosql_test

import (
	"github.com/cqlsy/leolib/leosql"
)

var (
	dbConf = &leosql.DbConfig{
		UserName: "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "leoSql_test",
	}
	testTable = leosql.Table{
		Name: "testTable",
		FieldList: map[string]leosql.DbType{
			"Id":   leosql.DbTypeIntAutoIncrement,
			"Name": leosql.DbTypeString,
			"Pwd":  leosql.DbTypeString,
		},
		PrimaryKey: "Id",
		UniqueKey: [][]string{
			[]string{"Id"},
		},
		NotNull: []string{"Name", "Pwd"},
	}
	tomlData = `
				[[testTable]]
				Id = "1"
				Name = "👮👮👮"
				Pwd = "111"
				[[testTable]]
				Id = "2"
				Name = "angelina2"
				Pwd = "222"
				[[testTable]]
				Id = "3"
				Name = "angelina3"
				Pwd = "333"
			`
)

func setTestTableData() {
	leosql.MustSetTableDataToml(tomlData)
}
