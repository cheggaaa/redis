package redis

import "errors"

var errZeroLengthValues = errors.New("Can't send zero values")

const (
	BITOP_AND = "AND"
	BITOP_OR = "OR"
	BITOP_XOR = "XOR"
	BITOP_NOT = "NOT"
)

func (r *Redis) Ping() (pong string, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("PING"),
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return 
	}
	pong = resp.line
	return
}

// Set the string value of a key
func (r *Redis) Set(key string, value []byte) (err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("SET"),
		[]byte(key),
		value,
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	if ! resp.isOk() {
		err = ErrNotOk
	}
	return
}

// Set the value and expiration of a key
func (r *Redis) SetEx(key string, seconds int, value []byte) (err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("SETEX"),
	}
	cmd.AddInt(seconds).Add(value)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	if ! resp.isOk() {
		err = ErrNotOk
	}
	return
}

// Get the value of a key
func (r *Redis) Get(key string) (bulk Bulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("GET"),
		[]byte(key),
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	bulk = resp.bulk
	return
}

// Append a value to a key
func (r *Redis) Append(key string, value []byte) (length int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("APPEND"),
		[]byte(key),
		value,
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	length = resp.integer
	return
}

// Decrement the integer value of a key by one
func (r *Redis) Decr(key string) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("DECR"),
		[]byte(key),
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}

// Decrement the integer value of a key by the given number
func (r *Redis) DecrBy(key string, number int) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("DECRBY"),
		[]byte(key),
	}
	cmd.AddInt(number)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}


// Increment the integer value of a key by one
func (r *Redis) Incr(key string) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("INCR"),
		[]byte(key),
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}

// Increment the integer value of a key by the given number
func (r *Redis) IncrBy(key string, number int) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("INCRBY"),
		[]byte(key),
	}
	cmd.AddInt(number)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}

// Set the string value of a key and return its old value
func (r *Redis) GetSet(key string, value []byte) (bulk Bulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("GETSET"),
		[]byte(key),
		value,
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	bulk = resp.bulk
	return
}

// Get a substring of the string stored at a key
func (r *Redis) GetRange(key string, start, end int) (bulk Bulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("GETRANGE"),
		[]byte(key),
	}
	cmd.AddInt(start).AddInt(end)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	bulk = resp.bulk
	return
}


// Set multiple keys to multiple values
func (r *Redis) MSet(values map[string][]byte) (err error) {
	if len(values) == 0 {
		return errZeroLengthValues
	}
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("MSET"),
	}
	
	for key, value := range values {
		cmd.AddString(key).Add(value)	
	}
	
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	if ! resp.isOk() {
		err = ErrNotOk
	}
	return
}

// Set multiple keys to multiple values, only if none of the keys exist
func (r *Redis) MSetNx(values map[string][]byte) (ok bool, err error) {
	if len(values) == 0 {
		return false, errZeroLengthValues
	}
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("MSETNX"),
	}
	
	for key, value := range values {
		cmd.AddString(key).Add(value)	
	}
	
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	if resp.integer > 0 {
		ok = true
	}
	return
}

// Get the values of all the given keys
func (r *Redis) MGet(keys... string) (values map[string]Bulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("MGET"),
	}
	
	for _, key := range keys {
		cmd.AddString(key)
	}
	
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	
	mb := resp.multiBulk	
	values = make(map[string]Bulk, len(keys))
	
	for i, key := range keys {
		values[key] = mb[i]
	}
	
	return
}

// Count set bits in a string
// For example:
//    Redis.BitCount("mykey", nil, nil) - will be BITCOUNT mykey
//    Redis.BitCount("mykey", &redis.NilInt{0}, &redis.Int{4}) - will be BITCOUNT mykey 0 4
//    Redis.BitCount("mykey", &redis.NilInt{2}, nil) - will be BITCOUNT mykey
func (r *Redis) BitCount(key string, start, end *NilInt) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("BITCOUNT"),
		[]byte(key),
	}
	if start != nil && end != nil {
		cmd.AddInt(start.I).AddInt(end.I)
	}	
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}

// Perform bitwise operations between strings
func (r *Redis) BitOp(destKey, op string, key... string) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("BITOP"),
		[]byte(destKey),
		[]byte(op),
	}
	for _, k := range key {
		cmd.AddString(k)
	}	
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}

// Returns the bit value at offset in the string value stored at key
func (r *Redis) GetBit(key string, offset int) (value int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("GETBIT"),
		[]byte(key),
	}
	cmd.AddInt(offset)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	value = resp.integer
	return
}

// Sets or clears the bit at offset in the string value stored at key
func (r *Redis) SetBit(key string, offset, value int) (result int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("SETBIT"),
		[]byte(key),
	}
	cmd.AddInt(offset).AddInt(value)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	result = resp.integer
	return
}