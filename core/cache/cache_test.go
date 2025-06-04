package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/baowk/dilu-core/config"
)

var memCache ICache
var redisCache ICache

type testCase struct {
	Key string `json:"key"`
	Val any    `json:"val"`
}

func (i testCase) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (tc *testCase) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tc)
}

type testCase2 struct {
	Key string `json:"key"`
	Val int    `json:"val"`
}

func (i testCase2) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (tc *testCase2) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, tc)
}

func TestRedisGS2(t *testing.T) {
	tc := testCase2{
		Key: "aaa",
		Val: 1,
	}
	//data, _ := json.Marshal(tc)
	key := "aaa"
	//fmt.Println("set", string(data))

	if err := redisCache.Set(key, tc, time.Minute*10); err != nil {
		t.Error(err)
	}
	if val, err := redisCache.Get(key); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("get:%v\n", val)

		var tc2 testCase2
		err := tc2.UnmarshalBinary([]byte(val))
		if err != nil {
			t.Error(err)
		}

		fmt.Printf("get:%+v\n", tc2)

		// var tc2 testCase2
		// //d2, _ := json.Marshal(val)
		// if err := json.Unmarshal([]byte(val), &tc2); err != nil {
		// 	t.Error(err)
		// }
		// fmt.Printf("%+v,%d,%d", tc2, tc.Val, tc2.Val)
		// if tc.Key != tc2.Key {
		// 	t.Error("redis get error")
		// }
	}
}

var m = make(map[string]int, 0)

var testGroup = []testCase{
	testCase{
		Key: "test1",
		Val: "test",
	},

	testCase{
		Key: "test2",
		Val: 1,
	},

	testCase{
		Key: "test3",
		Val: m,
	},
}

func init() {
	memCache = New(config.CacheCfg{Type: "memory"})

	redisCache = New(config.CacheCfg{
		Addr:       "127.0.0.1:6379",
		DB:         0,
		MasterName: "",
		Password:   "",
		Prefix:     "test:",
		Type:       "redis",
	})
	m["aaa"] = 1
	m["bbb"] = 2

}

func TestRedisGS(t *testing.T) {
	tc := testCase{
		Key: "aaa",
		Val: 1,
	}
	//data, _ := json.Marshal(tc)
	key := "aaa"
	//fmt.Println("set", string(data))

	if err := redisCache.Set(key, tc, time.Minute*10); err != nil {
		t.Error(err)
	}
	if val, err := redisCache.Get(key); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("get:%v\n", val)
		var tc2 testCase
		//d2, _ := json.Marshal(val)
		if err := json.Unmarshal([]byte(val), &tc2); err != nil {
			t.Error(err)
		}
		fmt.Printf("%+v,%d,%d", tc2, tc.Val, tc2.Val)
		if tc.Key != tc2.Key {
			t.Error("redis get error")
		}
	}
}

func TestGetRedis(t *testing.T) {
	c, err := GetRedisClient(redisCache)
	if err != nil {
		t.Error(err)
		return
	}
	if c == nil {
		t.Error("redis get error")
	}
	key := "aaa"

	str1, err1 := redisCache.Get(key)
	if err1 != nil {
		t.Error(err1)
	} else {
		fmt.Printf("get:%v\n", str1)
	}
	tc := testCase{
		Key: "aaa2",
		Val: 2,
	}
	data, _ := json.Marshal(tc)

	c.Set(context.Background(), redisCache.RealKey(key), data, 5*time.Minute)

	str, err := c.Get(context.Background(), redisCache.RealKey(key)).Result()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("get:%v\n", str)
	}

	str1, err1 = redisCache.Get(key)
	if err1 != nil {
		t.Error(err1)
	} else {
		fmt.Printf("get:%v\n", str1)
	}
}

func TestA(t *testing.T) {
	idx := 0
	memCache.Set(testGroup[idx].Key, testGroup[idx].Val, time.Duration(5)*time.Minute)
	str, err := memCache.Get(testGroup[idx].Key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	if str != testGroup[idx].Val {
		t.Errorf("The values of is not %v,%v \n", str, testGroup[idx].Val)
	}

}

func TestB(t *testing.T) {
	idx := 1
	memCache.Set(testGroup[idx].Key, testGroup[idx].Val, time.Duration(5)*time.Minute)
	str, err := memCache.Get(testGroup[idx].Key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	d, _ := strconv.Atoi(str)
	fmt.Printf("%v", d)
	if d != testGroup[idx].Val {
		t.Errorf("The values of is not %v,%v \n", d, testGroup[idx].Val)
	}

}

func TestC(t *testing.T) {
	idx := 2
	memCache.Set(testGroup[idx].Key, testGroup[idx].Val, time.Duration(5)*time.Minute)
	str, err := memCache.Get(testGroup[idx].Key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}

	//m := str.(map[string]int)
	d := make(map[string]int, 0)
	json.Unmarshal([]byte(str), &d)

	fmt.Printf("%v", m)

	// if d != testGroup[idx].Val {
	// 	t.Errorf("The values of is not %v,%v \n", d, testGroup[idx].Val)
	// }

}

func TestD(t *testing.T) {
	idx := 0

	if err := memCache.Set(testGroup[idx].Key, testGroup[idx].Val, time.Duration(5)*time.Minute); err != nil {
		t.Errorf("The values of is not %v\n", err)
	}

	if err := memCache.SetNX(testGroup[idx].Key, testGroup[idx].Val, time.Duration(5)*time.Minute); err != nil {
		t.Errorf("The values of is not %v\n", err)
	}
	str, err := memCache.Get(testGroup[idx].Key)
	if err != nil {
		t.Errorf("The values of is not %v\n", err)
	}

	fmt.Printf("res:%v", str)
}

func TestE(t *testing.T) {
	pairs := make(map[string]any, 0)
	keys := make([]string, 0)
	for _, v := range testGroup {
		pairs[v.Key] = v.Val
		keys = append(keys, v.Key)
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

	for i, v := range results {
		fmt.Printf("res[%s]:%v\n", keys[i], v)
	}

	m, err := GetMemoryClient(memCache)
	if err != nil {
		t.Errorf("Error getting memory client: %v", err)
		return
	}

	str, err := m.Get(testGroup[0].Key)
	if err != nil {
		t.Errorf("Error retrieving value: %v", err)
		return
	}

	fmt.Printf("res:%v\n", str)
}
