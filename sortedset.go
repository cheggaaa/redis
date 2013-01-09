package redis

// Add one or more members to a sorted set, or update its score if it already exists
func (r *Redis) ZAdd(key string, score int, member string) (result int, err error) {
	cmd := &Command{
		[]byte("ZADD"),
		[]byte(key),
	}
	cmd.AddInt(score).AddString(member)
	result, _, err = cmd.ExecuteInteger(r)
	return
}

// Get the number of members in a sorted set
func (r *Redis) ZCard(key string) (result int, err error) {
	cmd := &Command{
		[]byte("ZCARD"),
		[]byte(key),
	}
	result, _, err = cmd.ExecuteInteger(r)
	return
}

// Count the members in a sorted set with scores within the given values
func (r *Redis) ZCount(key string, min, max int) (result int, err error) {
	cmd := &Command{
		[]byte("ZCOUNT"),
		[]byte(key),
	}
	cmd.AddInt(min).AddInt(max)
	result, _, err = cmd.ExecuteInteger(r)
	return
}

// Increment the score of a member in a sorted set
func (r *Redis) ZIncrBy(key string, incr int, member string) (result int, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("ZINCRBY"),
		[]byte(key),
	}
	cmd.AddInt(incr).AddString(member)
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	result = resp.bulk.I()
	return
}

// Return a range of members in a sorted set, by index
func (r *Redis) ZRange(key string, start, stop int, withscores bool) (result MultiBulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("ZRANGE"),
		[]byte(key),
	}
	cmd.AddInt(start).AddInt(stop)
	if withscores {
		cmd.AddString("WITHSCORES")
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	result = resp.multiBulk
	return
}

// Return a range of members in a sorted set, by score (with limit)
func (r *Redis) ZRangeByScoreLimit(key string, min, max int, withscores bool, offset, count int) (result MultiBulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("ZRANGEBYSCORE"),
		[]byte(key),
	}
	cmd.AddInt(min).AddInt(max)
	if withscores {
		cmd.AddString("WITHSCORES")
	}
	if count > 0 {
		cmd.AddString("LIMIT").AddInt(offset).AddInt(count)	
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	result = resp.multiBulk
	return
}

// Return a range of members in a sorted set, by score
func (r *Redis) ZRangeByScore(key string, min, max int, withscores bool) (result MultiBulk, err error) {
	return r.ZRangeByScoreLimit(key, min, max, withscores, 0, 0)
}

// Determine the index of a member in a sorted set
func (r *Redis) ZRank(key, member string) (result int, exists bool, err error) {
	cmd := &Command{
		[]byte("ZRANK"),
		[]byte(key),
		[]byte(member),
	}
	result, exists, err = cmd.ExecuteInteger(r)
	return
}

// Remove one or more members from a sorted set
func (r *Redis) ZRem(key string, member ... string) (count int, err error) {
	cmd := &Command{
		[]byte("ZREM"),
		[]byte(key),	
	}
	for _, m := range member {
		cmd.AddString(m)
	}
	count, _, err = cmd.ExecuteInteger(r)
	return
}

// Remove all members in a sorted set within the given indexes
func (r *Redis) ZRemRangeByRank(key string, start, stop int) (result int, err error) {
	cmd := &Command{
		[]byte("ZREMRANGEBYRANK"),
		[]byte(key),
	}
	cmd.AddInt(start).AddInt(stop)
	result, _, err = cmd.ExecuteInteger(r)
	return
}

// Remove all members in a sorted set within the given scores
func (r *Redis) ZRemRangeByScore(key string, min, max int) (result int, err error) {
	cmd := &Command{
		[]byte("ZREMRANGEBYSCORE"),
		[]byte(key),
	}
	cmd.AddInt(min).AddInt(max)
	result, _, err = cmd.ExecuteInteger(r)
	return
}


// Return a range of members in a sorted set, by index, with scores ordered from high to low
func (r *Redis) ZRevRange(key string, start, stop int, withscores bool) (result MultiBulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("ZREVRANGE"),
		[]byte(key),
	}
	cmd.AddInt(start).AddInt(stop)
	if withscores {
		cmd.AddString("WITHSCORES")
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	result = resp.multiBulk
	return
}

// Return a range of members in a sorted set, by score, with scores ordered from high to low (with limit)
func (r *Redis) ZRevRangeByScoreLimit(key string, min, max int, withscores bool, offset, count int) (result MultiBulk, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("ZREVRANGEBYSCORE"),
		[]byte(key),
	}
	cmd.AddInt(min).AddInt(max)
	if withscores {
		cmd.AddString("WITHSCORES")
	}
	if count > 0 {
		cmd.AddString("LIMIT").AddInt(offset).AddInt(count)	
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	result = resp.multiBulk
	return
}

// Return a range of members in a sorted set, by score, with scores ordered from high to low
func (r *Redis) ZRevRangeByScore(key string, min, max int, withscores bool) (result MultiBulk, err error) {
	return r.ZRevRangeByScoreLimit(key, min, max, withscores, 0, 0)
}

// Determine the index of a member in a sorted set, with scores ordered from high to low
func (r *Redis) ZRevRank(key, member string) (result int, exists bool, err error) {
	cmd := &Command{
		[]byte("ZREVRANK"),
		[]byte(key),
		[]byte(member),
	}
	result, exists, err = cmd.ExecuteInteger(r)
	return
}

// Get the score associated with the given member in a sorted set
func (r *Redis) ZScore(key, member string) (result int, exists bool, err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("ZSCORE"),
		[]byte(key),
		[]byte(member),
	}
	err = r.Execute(cmd, resp)
	if err != nil {
		return
	}
	if resp.bulk != nil {
		exists = true
		result = resp.bulk.I()
	}
	return
}
