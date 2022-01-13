/**
 * Created by angelina on 2017/4/21.
 */

package leosql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/cqlsy/leolib/leosql/ast"
	"strings"
)

type leoTx struct {
	*sqlx.Tx
}

func BeginTx() (*leoTx, error) {
	tx, err := GetDb().Beginx()
	if err != nil {
		return nil, err
	}
	return &leoTx{tx}, nil
}

func (tx *leoTx) Commit() error {
	return tx.Tx.Commit()
}

func (tx *leoTx) Rollback() error {
	return tx.Tx.Rollback()
}

func (tx *leoTx) Query(query string, args ...interface{}) (output []map[string]string, err error) {
	rows, err := tx.Tx.Query(tx.Tx.Rebind(query), args...)
	colorPrint("[leosql.leoTx.Query]:Query=["+tx.Tx.Rebind(query)+"] args=", args)
	if err != nil {
		return nil, fmt.Errorf("[Query] leosql: [%s] data: [%s] err:[%s]",
			query, strings.Join(argsInterfaceToString(args), ","), err.Error())
	}
	defer rows.Close()
	// 获取table的字段数量
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	lenColumn := len(columns)
	for rows.Next() {
		rowArray := make([]interface{}, lenColumn)
		// 用*RawByte包装每一个值
		for k1 := range rowArray {
			var s sql.RawBytes
			rowArray[k1] = &s
		}
		// 读取这个row的全部值
		if err := rows.Scan(rowArray...); err != nil {
			return nil, err
		}
		rowMap := make(map[string]string)
		// 解包装，将全部字段取出来
		for rowIndex, rowName := range columns {
			rowMap[rowName] = string(*(rowArray[rowIndex].(*sql.RawBytes)))
		}
		// 这个row扔到output中
		output = append(output, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return
}

func (tx *leoTx) QueryOne(query string, args ...interface{}) (output map[string]string, err error) {
	list, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("no row find")
	}
	output = list[0]
	return output, err
}

func (tx *leoTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	colorPrint("[leosql.leoTx.Exec]:Query=["+tx.Tx.Rebind(query)+"] args=", args)
	return tx.Tx.Exec(tx.Tx.Rebind(query), args...)
}

func (tx *leoTx) Insert(tableName string, row map[string]string) (lastInsertID int, err error) {
	keyList := []string{}
	valueList := []string{}
	for key, value := range row {
		keyList = append(keyList, key)
		valueList = append(valueList, value)
	}
	keyStr := "`" + strings.Join(keyList, "`,`") + "`"
	valueStr := strings.Repeat("?,", len(row)-1) + "?"
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", tableName, keyStr, valueStr)
	result, err := tx.Exec(sql, argsStringToInterface(valueList...)...)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	lastInsertID = int(id)
	return
}

func (tx *leoTx) UpdateByID(tableName string, primaryKeyName string, row map[string]string) error {
	keyList := []string{}
	valueList := []string{}
	var primaryValue string
	for key, value := range row {
		if primaryKeyName == key {
			primaryValue = value
			continue
		}
		keyList = append(keyList, "`"+key+"`=?")
		valueList = append(valueList, value)
	}
	if primaryValue == "" {
		return fmt.Errorf("primaryKey %s not set", primaryKeyName)
	}
	valueList = append(valueList, primaryValue)
	updateStr := strings.Join(keyList, ",")
	// UPDATE User SET Name=?,Pwd=? WHERE id = ?
	sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE `%s` = ?", tableName, updateStr, primaryKeyName)
	_, err := tx.Exec(sql, argsStringToInterface(valueList...)...)
	if err != nil {
		return err
	}
	return nil
}

func (tx *leoTx) DeleteByID(tableName, fieldName, value string) error {
	deleteSql := fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = ?", tableName, fieldName)
	_, err := tx.Exec(deleteSql, value)
	return err
}

func (tx *leoTx) GetOneWhere(tableName, fieldName, value string) (map[string]string, error) {
	getSql := fmt.Sprintf("SELECT * FROM `%s` WHERE `%s` = ?", tableName, fieldName)
	return tx.QueryOne(getSql, value)
}

func (tx *leoTx) GetAllInTable(tableName string) ([]map[string]string, error) {
	getSql := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	return tx.Query(getSql)
}

func (tx *leoTx) RunSelectCommand(selectCommand *ast.SelectCommand) (mapValue []map[string]string, err error) {
	prepareSql, parameterList := selectCommand.GetPrepareParameter()
	mapValue, err = tx.Query(prepareSql, argsStringToInterface(parameterList...)...)
	return
}

func (tx *leoTx) IsExist(tableName string, row map[string]string) bool {
	where := ast.NewAndWhereCondition()
	for k, v := range row {
		where = where.AddPrepare(fmt.Sprintf("%s=?", k), v)
	}
	selectCommand := ast.NewSelectCommand().From(tableName).WhereObj(where)
	info, err := tx.RunSelectCommand(selectCommand)
	if err != nil || len(info) == 0 {
		return false
	}
	return true
}
