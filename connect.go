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
	var written, n int
	cmdB := cmd.ProtoRequest()
	for written < len(cmdB) {
		n, err = c.writer.Write(cmdB[written:])
		if err != nil {
			return
		}
		err = c.writer.Flush()
		if err != nil {
			return
		}
		written += n
	}
	err = r.Read(c.reader)
	return	
}