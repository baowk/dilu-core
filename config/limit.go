package config

import "time"

type AccessLimit struct {
	Enable   bool          `mapstructure:"enable" json:"enable" yaml:"enable"`       //是否启用
	Duration time.Duration `mapstructure:"duration" json:"duration" yaml:"duration"` //时长周期
	Total    int           `mapstructure:"total" json:"total" yaml:"total"`          //周期内最大访问次数，超过拒绝
}

func (s *AccessLimit) GetDuration() time.Duration {
	if s.Duration < 1 {
		s.Duration = time.Second
	}
	return s.Duration
}

func (s *AccessLimit) GetTotal() int {
	if s.Total < 0 {
		s.Total = 100
	}
	return s.Total
}
