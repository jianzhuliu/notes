package consistenthash

import (
	"reflect"
	"strconv"
	"testing"
)

//hash算法，转换成数字,方便测试
var hashFuncTest = func(data []byte) uint32 {
	i, _ := strconv.Atoi(string(data))
	return uint32(i)
}

func TestHashObj(t *testing.T) {
	//复制3个虚拟节点
	hashObj := New(3, hashFuncTest)

	//添加真实节点
	hashObj.Add("5", "3", "7")
	//t.Log(hashObj.hashValues)
	//3 5 7 13 15 17 23 25 27

	expect := []int{3, 5, 7, 13, 15, 17, 23, 25, 27}
	if !reflect.DeepEqual(hashObj.hashValues, expect) {
		t.Fatalf("hashValues not match, expect %v, but %v got", expect, hashObj.hashValues)
	}

	testCases := map[string]string{
		"3":  "3",
		"8":  "3",
		"14": "5",
		"26": "7",
		"42": "3",
	}

	for k, v := range testCases {
		if got := hashObj.Get(k); got != v {
			t.Fatalf("looking for %s, expect %s, but %s got", k, v, got)
		}
	}

	//新增节点
	hashObj.Add("9")
	testCases["8"] = "9"

	for k, v := range testCases {
		if got := hashObj.Get(k); got != v {
			t.Fatalf("after add new note,looking for %s, expect %s, but %s got", k, v, got)
		}
	}
}
