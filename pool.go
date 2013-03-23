package redis

import ()

func newPool(c *Config) (pool *ConnPool, err error) {
	pool = &ConnPool{
		make(chan *Command),
		make([]*Conn, c.MaxConnections),
	}
	for i := 0; i < c.MaxConnections; i++ {
		pool.conns[i], err = newConn("tcp", c.Addr, c.ConnectTimeout)
		if err != nil {
			return
		}
	}
	return
}

type ConnPool struct {
	in    chan *Command
	conns []*Conn
}
