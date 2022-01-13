package fileUtil

import (
	"testing"
)

func TestName(t *testing.T) {
	println(removePreAndSuf("/dsdas", "/"))
	println(removePreAndSuf("/dsdas/dasd/ds", "/"))
	println(removePreAndSuf("/dsdas/", "/"))
	println(removePreAndSuf("dsdas/", "/"))
}
