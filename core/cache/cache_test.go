package cache

import (
	"fmt"
	"testing"
	"time"
)

var memCache ICache
var redisCache ICache

func init() {
	memCache = NewMemory()
	m["aaa"] = 1
	m["bbb"] = 2
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
	d, _ := str.(int)
	fmt.Printf("%v", d)
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

	m := str.(map[string]int)
	// d := make(map[string]int, 0)
	// json.Unmarshal([]byte(str), &d)

	fmt.Printf("%v", m)

	// if d != testGroup[idx].val {
	// 	t.Errorf("The values of is not %v,%v \n", d, testGroup[idx].val)
	// }

}

func TestD(t *testing.T) {
	idx := 0

	if err := memCache.Set(testGroup[idx].key, testGroup[idx].val, time.Duration(5)*time.Minute); err != nil {
		t.Errorf("The values of is not %v\n", err)
	}

	if err := memCache.SetNX(testGroup[idx].key, testGroup[idx].val, time.Duration(5)*time.Minute); err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	str, err := memCache.Get(testGroup[idx].key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}

	fmt.Printf("res:%v", str)
}

func TestE(t *testing.T) {
	pairs := make(map[string]any, 0)
	keys := make([]string, 0)
	for _, v := range testGroup {
		pairs[v.key] = v.val
		keys = append(keys, v.key)
	}

	if err := memCache.MSet(pairs); err != nil {
		t.Errorf("Error setting multiple values: %v", err)
		return
	}

	results, err := memCache.MGet(keys...)
	if err != nil {
		t.Errorf("Error retrieving multiple values: %v", err)
		return
	}

	for i, v := range results.([]interface{}) {
		fmt.Printf("res[%s]:%v\n", keys[i], v)
	}

	m, err := GetMemoryClient(memCache)
	if err != nil {
		t.Errorf("Error getting memory client: %v", err)
		return
	}

	str, err := m.Get(testGroup[0].key)
	if err != nil {
		t.Errorf("Error retrieving value: %v", err)
		return
	}

	fmt.Printf("res:%v\n", str)
}
