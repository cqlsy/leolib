/**
 * Created by angelina on 2017/4/16.
 */

package leosql_test

import (
	"github.com/cqlsy/leolib/leosql"
	"testing"
)

func TestMustSetTableDataToml(t *testing.T) {
	initTestDbTable()
	leosql.MustSetTableDataToml(tomlData)
}
