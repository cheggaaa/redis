package redis

import (
	"errors"
	"bufio"
	//"fmt"
	"strconv"
)

type Reader interface {
	Read(r *bufio.Reader) (err error)
}

type ReaderBase struct {
	rType string	
	line string
	bulk Bulk
	multiBulk MultiBulk
	err string
	integer int
}


const (
	// types
	rt_line = "+"
	rt_integer = ":"
	rt_error = "-"
	rt_bulk = "$"
	rt_multi_bulk = "*"
)

func (rb *ReaderBase) Read(r *bufio.Reader) (err error) {
	// detect type
	err = rb.detectType(r)
	if err != nil {
		return
	}
	switch rb.rType {
		case rt_line:
			err = rb.readLine(r)
			break
		case rt_bulk:
			err = rb.readBulk(r)
			break
		case rt_integer:
			err = rb.readIneger(r)
			break
		case rt_error:
			err = rb.readError(r)
			break
		case rt_multi_bulk:
			err = rb.readMultiBulk(r)
	}
	return
}


func (rb *ReaderBase) Error() (err error) {
	if rb.rType == rt_error {
		err = errors.New(rb.err)
	}
	return
}

func (rb *ReaderBase) isNil() bool {
	return rb.rType == rt_bulk && rb.bulk == nil 
}

func (rb *ReaderBase) isOk() bool {
	return rb.rType == rt_line && rb.line == "OK"
}


func (rb *ReaderBase) detectType(r *bufio.Reader) (err error) {
	firstByte, err := r.ReadByte()
	if err != nil {
		return
	}
	//fmt.Println(string(firstByte))
	switch string(firstByte) {
		case rt_line, rt_integer, rt_error, rt_bulk, rt_multi_bulk:
			rb.rType = string(firstByte)
			break
		default:
		err = errors.New("Unknown server response type - " + string(firstByte))
	}
	return
}


func (rb *ReaderBase) readLine(r *bufio.Reader) (err error) {
	b, err := rb.readToCrLf(r)
	if err == nil {
		rb.line = string(b)
	}
	return	
}

func (rb *ReaderBase) readIneger(r *bufio.Reader) (err error) {
	lb, err := rb.readToCrLf(r)
	if err != nil {
		return
	}
	if len(lb) > 0  && string(lb[0]) == rt_bulk {
		lb = lb[1:]
	}
	rb.integer, err = strconv.Atoi(string(lb))
	return	
}

func (rb *ReaderBase) readError(r *bufio.Reader) (err error) {
	b, err := rb.readToCrLf(r)
	if err == nil {
		rb.err = string(b)
	}
	return	
}


func (rb *ReaderBase) readBulk(r *bufio.Reader) (err error) {
	err = rb.readIneger(r)
	if err != nil {
		return
	}
	if rb.integer == -1 {
		rb.bulk = nil
		return
	}
	rb.bulk = make([]byte, rb.integer)
	
	n := 0
	written := 0
	
	for written < rb.integer {
		buf := make([]byte, rb.integer - written)
		n, err = r.Read(buf)
		if err != nil {
			return
		}
		
		if n > 0 {
			copy(rb.bulk[written:], buf[0:n])
			written += n
		}
	}
	err = rb.readCrLf(r)
	return	
}


func (rb *ReaderBase) readMultiBulk(r *bufio.Reader) (err error) {
	err = rb.readIneger(r)
	if err != nil {
		return
	}
	count := rb.integer
	rb.multiBulk = make(MultiBulk, count)
	for i := 0; i < count; i++ {
		
		err = rb.readBulk(r)
		if err != nil {
			return
		}
		rb.multiBulk[i] = rb.bulk
	}
	return
}

func (rb *ReaderBase) readToCrLf(r *bufio.Reader) (b []byte, err error) {
	b, err = r.ReadBytes('\n')
	if len(b) < 2 {
		return nil, errors.New("Unexpected response")
	}
	return b[0:len(b)-2], nil
}

func (rb *ReaderBase) readCrLf(r *bufio.Reader) (err error) {
	_, err = r.ReadBytes('\n')
	return
}
