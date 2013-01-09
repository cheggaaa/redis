package redis

import (
	"net"
	"sync"
	"time"
	"bufio"
)

func newConn(netw, addr string, timeout time.Duration) (conn *Conn, err error) {
	c, err := net.DialTimeout(netw, addr, timeout)
	if err != nil {
		return
	}
	conn = &Conn{
		reader : bufio.NewReader(c), 
		writer : bufio.NewWriter(c), 
		c      : c,
		m      : &sync.Mutex{},
	}
	return
}

type Conn struct {
	reader *bufio.Reader
	writer *bufio.Writer
	c       net.Conn
	m       *sync.Mutex
}

func (c *Conn) WriteRead(cmd *Command, r Reader) (err error) {
	c.m.Lock()
	defer c.m.Unlock()
	_, err = c.writer.Write(cmd.ProtoRequest())
	if err != nil {
		return
	}
	err = c.writer.Flush()
	if err != nil {
		return
	}
	err = r.Read(c.reader)
	return	
}