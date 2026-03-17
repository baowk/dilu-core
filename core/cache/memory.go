package cache

import (
	"encoding/json"
	"errors"
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
}

func (*Memory) Type() string {
	return "memory"
}

func (m *Memory) IsRedis() bool {
	return false
}

func (m *Memory) RealKey(key string) string {
	return key
}

func (m *Memory) Get(key string) (string, error) {
	item, err := m.getItem(key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, nil
}

func (m *Memory) getItem(key string) (*item, error) {
	i, ok := m.items.Load(key)
	if !ok {
		return nil, errors.New("not exist")
	}
	switch v := i.(type) {
	case *item:
		if v.Expired.Before(time.Now()) {
			m.items.Delete(key)
			return nil, errors.New("not exist")
		}
		return v, nil
	default:
		return nil, fmt.Errorf("value of %s type error", key)
	}
}

func (m *Memory) Set(key string, val interface{}, expiration time.Duration) error {
	s, err := cast.ToStringE(val)
	if err != nil {
		bs, err := json.Marshal(val)
		if err != nil {
			return err
		}
		s = string(bs)
	}
	it := &item{
		Value:   s,
		Expired: time.Now().Add(expiration),
	}
	return m.setItem(key, it)
}

func (m *Memory) SetNX(key string, val interface{}, expiration time.Duration) error {
	if m.Exists(key) {
		return errors.New("key exist")
	}
	return m.Set(key, val, expiration)
}

func (m *Memory) setItem(key string, it *item) error {
	m.items.Store(key, it)
	return nil
}

func (m *Memory) Del(key string) error {
	m.items.Delete(key)
	return nil
}

func (m *Memory) HGet(hk, key string) (any, error) {
	it, err := m.getItem(hk + key)
	if err != nil || it == nil {
		return "", err
	}
	return it.Value, err
}

func (m *Memory) HDel(hk, key string) error {
	m.items.Delete(hk + key)
	return nil
}

func (m *Memory) Incr(key string) (int64, error) {
	return m.calculate(key, 1)
}

func (m *Memory) Decr(key string) (int64, error) {
	return m.calculate(key, -1)
}

// calculate 使用 CompareAndSwap 模式实现原子递增/递减
func (m *Memory) calculate(key string, num int64) (int64, error) {
	for {
		old, ok := m.items.Load(key)
		if !ok {
			return 0, fmt.Errorf("%s not exist", key)
		}
		it, ok := old.(*item)
		if !ok {
			return 0, fmt.Errorf("value of %s type error", key)
		}
		if it.Expired.Before(time.Now()) {
			m.items.Delete(key)
			return 0, errors.New("not exist")
		}

		n, err := cast.ToInt64E(it.Value)
		if err != nil {
			return 0, err
		}
		n += num
		newItem := &item{
			Value:   strconv.FormatInt(n, 10),
			Expired: it.Expired,
		}
		if m.items.CompareAndSwap(key, old, newItem) {
			return n, nil
		}
		// CAS 失败，重试
	}
}

func (m *Memory) Expire(key string, dur time.Duration) error {
	old, ok := m.items.Load(key)
	if !ok {
		return fmt.Errorf("%s not exist", key)
	}
	it, ok := old.(*item)
	if !ok {
		return fmt.Errorf("%s type error", key)
	}
	newItem := &item{
		Value:   it.Value,
		Expired: time.Now().Add(dur),
	}
	m.items.CompareAndSwap(key, old, newItem)
	return nil
}

func (m *Memory) ExpireAt(key string, tm time.Time) error {
	old, ok := m.items.Load(key)
	if !ok {
		return fmt.Errorf("%s not exist", key)
	}
	it, ok := old.(*item)
	if !ok {
		return fmt.Errorf("%s type error", key)
	}
	newItem := &item{
		Value:   it.Value,
		Expired: tm,
	}
	m.items.CompareAndSwap(key, old, newItem)
	return nil
}

func (m *Memory) Exists(key string) bool {
	_, err := m.getItem(key)
	return err == nil
}

func (m *Memory) MGet(keys ...string) ([]any, error) {
	var values []any
	for _, key := range keys {
		it, err := m.getItem(key)
		if err != nil {
			return nil, err
		}
		if it == nil {
			return nil, fmt.Errorf("%s not exist", key)
		}
		values = append(values, it.Value)
	}
	return values, nil
}

func (m *Memory) MSet(pairs map[string]any) error {
	for key, v := range pairs {
		m.Set(key, v, time.Hour*24*365)
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
