package redis

// Delete a key
func (r *Redis) Del(key ...string) (count int, err error) {
	cmd := &Command{[]byte("DEL")}
	for _, k := range key {
		cmd.AddString(k)
	}
	count, _, err = cmd.ExecuteInteger(r)
	return
}

// Determine if a key exists
func (r *Redis) Exists(key string) (exists bool, err error) {
	cmd := &Command{
		[]byte("EXISTS"),
		[]byte(key),
	}
	i, _, err := cmd.ExecuteInteger(r)
	if err == nil {
		exists = i > 0
	}
	return
}

// Set a key's time to live in seconds
func (r *Redis) Expire(key string, seconds int) (ok bool, err error) {
	cmd := &Command{
		[]byte("EXPIRE"),
		[]byte(key),
	}
	cmd.AddInt(seconds)
	i, _, err := cmd.ExecuteInteger(r)
	if err == nil {
		ok = i > 0
	}
	return
}

// Set the expiration for a key as a UNIX timestamp
func (r *Redis) ExpireAt(key string, timestamp int) (ok bool, err error) {
	cmd := &Command{
		[]byte("EXPIREAT"),
		[]byte(key),
	}
	cmd.AddInt(timestamp)
	i, _, err := cmd.ExecuteInteger(r)
	if err == nil {
		ok = i > 0
	}
	return
}

// Set a key's time to live in milliseconds
func (r *Redis) PExpire(key string, milliseconds int) (ok bool, err error) {
	cmd := &Command{
		[]byte("PEXPIRE"),
		[]byte(key),
	}
	cmd.AddInt(milliseconds)
	i, _, err := cmd.ExecuteInteger(r)
	if err == nil {
		ok = i > 0
	}
	return
}

// Set the expiration for a key as a UNIX timestamp specified in millisecond
func (r *Redis) PExpireAt(key string, timestamp int) (ok bool, err error) {
	cmd := &Command{
		[]byte("PEXPIREAT"),
		[]byte(key),
	}
	cmd.AddInt(timestamp)
	if i, _, err := cmd.ExecuteInteger(r); err == nil {
		ok = i > 0
	}
	return
}

// Get the time to live for a key
func (r *Redis) TTL(key string) (ttl int, err error) {
	cmd := &Command{
		[]byte("TTL"),
		[]byte(key),
	}
	ttl, _, err = cmd.ExecuteInteger(r)
	return
}

// Get the time to live for a key in milliseconds
func (r *Redis) PTTL(key string) (ttl int, err error) {
	cmd := &Command{
		[]byte("PTTL"),
		[]byte(key),
	}
	ttl, _, err = cmd.ExecuteInteger(r)
	return
}

// Find all keys matching the given pattern
func (r *Redis) Keys(pattern string) (keys MultiBulk, err error) {
	cmd := &Command{
		[]byte("KEYS"),
		[]byte(pattern),
	}
	resp, err := cmd.Execute(r)
	if err != nil {
		return
	}
	keys = resp.multiBulk
	return
}

// Return a serialized version of the value stored at the specified key
func (r *Redis) Dump(key string) (result []byte, err error) {
	cmd := &Command{
		[]byte("DUMP"),
		[]byte(key),
	}
	if resp, err := cmd.Execute(r); err == nil {
		result = []byte(resp.bulk)
	}
	return
}

// Move a key to another database
func (r *Redis) Move(key string, db int) (moved bool, err error) {
	cmd := &Command{
		[]byte("MOVE"),
		[]byte(key),
	}
	cmd.AddInt(db)
	if i, _, err := cmd.ExecuteInteger(r); err == nil && i > 0 {
		moved = true
	}
	return
}

// Inspect the internals of Redis objects
func (r *Redis) Object(subcommand string, args ...string) (intValue int, bulkValue Bulk, err error) {
	cmd := &Command{
		[]byte("OBJECT"),
		[]byte(subcommand),
	}
	for _, arg := range args {
		cmd.AddString(arg)
	}
	resp, err := cmd.Execute(r)
	if err != nil {
		return
	}
	switch resp.rType {
	case rt_bulk:
		bulkValue = resp.bulk
	case rt_integer:
		intValue = resp.integer
	}
	return
}

// Remove the expiration from a key
func (r *Redis) Persist(key string) (timeoutRemoved bool, err error) {
	cmd := &Command{
		[]byte("PERSIST"),
		[]byte(key),
	}

	if i, _, err := cmd.ExecuteInteger(r); err == nil && i > 0 {
		timeoutRemoved = true
	}
	return
}

// Return a random key from the currently selected database
func (r *Redis) RandomKey() (key Bulk, err error) {
	cmd := &Command{
		[]byte("RANDOMKEY"),
	}
	if resp, err := cmd.Execute(r); err == nil {
		key = resp.bulk
	}
	return
}

// Rename a key
func (r *Redis) Rename(key, newKey string) (err error) {
	cmd := &Command{
		[]byte("RENAME"),
		[]byte(key),
		[]byte(newKey),
	}
	resp, err := cmd.Execute(r)
	if err != nil {
		return
	}
	if !resp.isOk() {
		err = ErrNotOk
	}
	return
}

// Rename a key, only if the new key does not exist
func (r *Redis) RenameNx(key, newKey string) (renamed bool, err error) {
	cmd := &Command{
		[]byte("RENAMENX"),
		[]byte(key),
		[]byte(newKey),
	}

	if i, _, err := cmd.ExecuteInteger(r); err == nil && i > 0 {
		renamed = true
	}
	return
}

// Create a key using the provided serialized value, previously obtained using DUMP
func (r *Redis) Restore(key string, ttl int, value []byte) (err error) {
	cmd := &Command{
		[]byte("RESTORE"),
		[]byte(key),
	}
	cmd.AddInt(ttl).Add(value)
	resp, err := cmd.Execute(r)
	if err != nil {
		return
	}
	if !resp.isOk() {
		err = ErrNotOk
	}
	return
}

// Determine the type stored at key
func (r *Redis) Type(key string) (keyType string, err error) {
	cmd := &Command{
		[]byte("TYPE"),
		[]byte(key),
	}
	resp, err := cmd.Execute(r)
	if err != nil {
		return
	}
	keyType = resp.line
	return
}
