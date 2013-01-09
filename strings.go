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
