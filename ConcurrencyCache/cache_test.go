
package ConcurrencyCache

import (
	"testing"
	"reflect"
)

//用一个 map 模拟耗时的数据库。
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T){

}