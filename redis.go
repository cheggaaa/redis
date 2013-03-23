package redis

import (
	"errors"
)

var (
	ErrNotOk = errors.New("Response not OK")
)

const (
	// types
	TYPE_STRING = "string"
	TYPE_LIST   = "list"
	TYPE_SET    = "set"
	TYPE_ZSET   = "zset"
	TYPE_HASH   = "hash"
	TYPE_NONE   = "none"
)

type Redis struct {
	Config *Config
	conn   *Conn
}

func (r *Redis) connect() (err error) {
	r.conn, err = newConn("tcp", r.Config.Addr, r.Config.ConnectTimeout)
	return
}

func (r *Redis) Execute(c *Command, response Reader) (err error) {
	return r.conn.WriteRead(c, response)
}
