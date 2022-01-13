/**
 * Created by angelina on 2017/4/15.
 */

package leosql_test

import (
	"github.com/cqlsy/leolib/leosql"
	"testing"
)

func initTestDbTable() {
	leosql.MustSetDbConfig(dbConf)
	leosql.InitDbWithoutDbName()
	leosql.MustCreateDb()
	leosql.InitDb()
	leosql.MustCreateTable(testTable)
}

func TestMustCreateDb(t *testing.T) {
	leosql.MustSetDbConfig(dbConf)
	leosql.InitDbWithoutDbName()
	leosql.MustCreateDb()
}
