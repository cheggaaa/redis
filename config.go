package redis

import (
	"time"
)

type Config struct {
	Addr string
	MaxConnections int
	ConnectTimeout time.Duration
}

func (c *Config) Init() error {
	if c.Addr == "" {
		c.Addr = ":6379"
	}
	if c.MaxConnections <= 0 {
		c.MaxConnections = 10
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = time.Second * 5
	}
	return nil
}