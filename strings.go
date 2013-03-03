package redis

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