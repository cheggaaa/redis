package redis

func (r *Redis) Select(db int) (err error) {
	resp := &ReaderBase{}
	cmd := &Command{
		[]byte("SELECT"),
	}
	cmd.AddInt(db)
	if err = r.Execute(cmd, resp); err != nil {
		return
	}
	if !resp.isOk() {
		err = ErrNotOk
	}
	return
}
