package cache

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type item struct {
	Value   any
	Expired time.Time
}

// NewMemory memory模式
func NewMemory() *Memory {
	return &Memory{
		items: new(sync.Map),
	}
}

type Memory struct {
	items *sync.Map
	mutex sync.RWMutex
}

func (*Memory) Type() string {
	return "memory"
}

func (m *Memory) Get(key string) (any, error) {
	item, err := m.getItem(key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, nil
}

func (m *Memory) getItem(key string) (*item, error) {
	var err error
	i, ok := m.items.Load(key)
	if !ok {
		return nil, errors.New("not exist")
	}
	switch i.(type) {
	case *item:
		item := i.(*item)
		if item.Expired.Before(time.Now()) {
			//过期
			_ = m.del(key)
			//过期后删除
			return nil, errors.New("not exist")
		}
		return item, nil
	default:
		err = fmt.Errorf("value of %s type error", key)
		return nil, err
	}
}

func (m *Memory) Set(key string, val interface{}, expiration time.Duration) error {
	// s, err := cast.ToStringE(val)
	// if err != nil {
	// 	bs, err := json.Marshal(val)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return err
	// 	}
	// 	s = string(bs)
	// }
	item := &item{
		Value:   val,
		Expired: time.Now().Add(expiration),
	}
	return m.setItem(key, item)
}

func (m *Memory) SetNX(key string, val interface{}, expiration time.Duration) error {
	if m.Exists(key) {
		return errors.New("key exist")
	}
	return m.Set(key, val, expiration)
}

func (m *Memory) setItem(key string, item *item) error {
	m.items.Store(key, item)
	return nil
}

func (m *Memory) Del(key string) error {
	return m.del(key)
}

func (m *Memory) del(key string) error {
	m.items.Delete(key)
	return nil
}

func (m *Memory) HGet(hk, key string) (any, error) {
	item, err := m.getItem(hk + key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, err
}

func (m *Memory) HDel(hk, key string) error {
	return m.del(hk + key)
}

func (m *Memory) Incr(key string) (int64, error) {
	return m.calculate(key, 1)
}

func (m *Memory) Decr(key string) (int64, error) {
	return m.calculate(key, -1)
}

func (m *Memory) calculate(key string, num int64) (int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return 0, err
	}

	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return 0, err
	}
	var n int64
	n, err = cast.ToInt64E(item.Value)
	if err != nil {
		return 0, err
	}
	n += num
	item.Value = strconv.FormatInt(n, 10)
	return n, m.setItem(key, item)
}

func (m *Memory) Expire(key string, dur time.Duration) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}
	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	item.Expired = time.Now().Add(dur)
	return m.setItem(key, item)
}

func (m *Memory) Exists(key string) bool {
	_, err := m.getItem(key)
	return err != nil
}

func (m *Memory) MGet(keys ...string) (any, error) {
	var values []any
	for _, key := range keys {
		item, err := m.getItem(key)
		if err != nil {
			return nil, err
		}
		if item == nil {
			err = fmt.Errorf("%s not exist", key)
			return nil, err
		}
		values = append(values, item.Value)
	}
	return values, nil
}

func (m *Memory) MSet(pairs map[string]any) error {
	for key, v := range pairs {
		item := &item{
			Value:   v, // 直接存储 value，不进行类型断言
			Expired: time.Now().Add(time.Hour * 24 * 365),
		}
		m.items.Store(key, item)
	}
	return nil
}

func (m *Memory) GetClient() *Memory {
	return m
}

func GetMemoryClient(c ICache) (*Memory, error) {
	if c != nil && c.Type() == "memory" {
		return c.(*Memory), nil
	}
	return nil, errors.ErrUnsupported
}
