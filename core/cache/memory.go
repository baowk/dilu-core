package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type item struct {
	Value   string
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

func (m *Memory) Get(key string) (string, error) {
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
		return nil, nil
	}
	switch i.(type) {
	case *item:
		item := i.(*item)
		if item.Expired.Before(time.Now()) {
			//过期
			_ = m.del(key)
			//过期后删除
			return nil, nil
		}
		return item, nil
	default:
		err = fmt.Errorf("value of %s type error", key)
		return nil, err
	}
}

func (m *Memory) Set(key string, val interface{}, expiration time.Duration) error {
	s, err := cast.ToStringE(val)
	if err != nil {
		bs, err := json.Marshal(val)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		s = string(bs)
	}
	item := &item{
		Value:   s,
		Expired: time.Now().Add(expiration),
	}
	return m.setItem(key, item)
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

func (m *Memory) HGet(hk, key string) (string, error) {
	item, err := m.getItem(hk + key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, err
}

func (m *Memory) HDel(hk, key string) error {
	return m.del(hk + key)
}

func (m *Memory) Incr(key string) error {
	return m.calculate(key, 1)
}

func (m *Memory) Decr(key string) error {
	return m.calculate(key, -1)
}

func (m *Memory) calculate(key string, num int) error {
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
	var n int
	n, err = cast.ToIntE(item.Value)
	if err != nil {
		return err
	}
	n += num
	item.Value = strconv.Itoa(n)
	return m.setItem(key, item)
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

func (m *Memory) GetClient() *Memory {
	return m
}
