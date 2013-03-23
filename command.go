package redis

import (
	"strconv"
)

const (
	crlf    = "\r\n"
	numargs = "*"
	arglen  = "$"
)

type Command [][]byte

func (c *Command) ProtoRequest() (request []byte) {
	request = []byte(numargs + strconv.Itoa(len(*c)) + crlf)
	for _, arg := range *c {
		request = append(request, []byte(arglen+strconv.Itoa(len(arg))+crlf)...)
		request = append(request, arg...)
		request = append(request, []byte(crlf)...)
	}
	return
}

func (c *Command) AddString(s string) *Command {
	*c = append(*c, []byte(s))
	return c
}

func (c *Command) AddInt(i int) *Command {
	return c.AddIntB(i, 10)
}

func (c *Command) AddIntB(i, b int) *Command {
	*c = append(*c, strconv.AppendInt(nil, int64(i), b))
	return c
}

func (c *Command) AddBool(b bool) *Command {
	*c = append(*c, strconv.AppendBool(nil, b))
	return c
}

func (c *Command) Add(b []byte) *Command {
	*c = append(*c, b)
	return c
}

func (c *Command) String() (s string) {
	for i, arg := range *c {
		if i > 0 {
			s += " "
		}
		s += string(arg)
	}
	return
}

func (c *Command) Execute(r *Redis) (reader *ReaderBase, err error) {
	reader = &ReaderBase{}
	if err = r.Execute(c, reader); err != nil {
		return
	}
	err = reader.Error()
	return
}

func (c *Command) ExecuteInteger(r *Redis) (i int, isNil bool, err error) {
	reader, err := c.Execute(r)
	if reader.rType == rt_integer {
		i = reader.integer
		return
	}
	isNil = reader.isNil()
	return
}
