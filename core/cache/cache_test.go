package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var memCache ICache
var redisCache ICache

func init() {
	memCache = NewMemory()
	m["aaa"] = 1
}

type testCase struct {
	key string
	val any
}

var m = make(map[string]int, 0)

var testGroup = []testCase{
	testCase{
		key: "test1",
		val: "test",
	},

	testCase{
		key: "test2",
		val: 1,
	},

	testCase{
		key: "test3",
		val: m,
	},
}

func TestA(t *testing.T) {
	idx := 0
	memCache.Set(testGroup[idx].key, testGroup[idx].val, time.Duration(5)*time.Minute)
	str, err := memCache.Get(testGroup[idx].key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	if str != testGroup[idx].val {
		t.Errorf("The values of is not %v,%v \n", str, testGroup[idx].val)
	}

}

func TestB(t *testing.T) {
	idx := 1
	memCache.Set(testGroup[idx].key, testGroup[idx].val, time.Duration(5)*time.Minute)
	str, err := memCache.Get(testGroup[idx].key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	d, _ := strconv.Atoi(str)
	if d != testGroup[idx].val {
		t.Errorf("The values of is not %v,%v \n", d, testGroup[idx].val)
	}

}

func TestC(t *testing.T) {
	idx := 2
	memCache.Set(testGroup[idx].key, testGroup[idx].val, time.Duration(5)*time.Minute)
	str, err := memCache.Get(testGroup[idx].key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	d := make(map[string]int, 0)
	json.Unmarshal([]byte(str), &d)

	fmt.Printf("%v", d)

	// if d != testGroup[idx].val {
	// 	t.Errorf("The values of is not %v,%v \n", d, testGroup[idx].val)
	// }

}
