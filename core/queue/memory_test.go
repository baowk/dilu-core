package queue

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestMemory_Append(t *testing.T) {
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
	}
	type args struct {
		name    string
		message IMessager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test01",
			fields{},
			args{
				name: "test",
				message: &Message{BaseMsg{
					ID:     "",
					Stream: "test",
					Values: map[string]interface{}{
						"key": "value",
					},
				}, 3, sync.RWMutex{},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(100)
			if err := m.Append(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Register(t *testing.T) {
	log.SetFlags(19)
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
	}
	type args struct {
		name string
		f    ConsumerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test01",
			fields{},
			args{
				name: "test",
				f: func(message IMessager) error {
					fmt.Println(message.GetValues())
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(100)
			m.Register(tt.name, tt.args.f)
			if err := m.Append(&Message{BaseMsg{
				Stream: "test",
				Values: map[string]interface{}{
					"key": "value",
				},
			}, 3, sync.RWMutex{}}); err != nil {
				t.Error(err)
				return
			}
			go func() {
				m.Run()
			}()
			time.Sleep(3 * time.Second)
			m.Shutdown()
		})
	}
}
