package redis

func (r *Redis) Select(db int) (err error) {
	cmd := &Command{
		[]byte("SELECT"),
	}	
	resp, err := cmd.AddInt(db).Execute(r)
	if err != nil {
		return
	}
	if ! resp.isOk() {
		return ErrNotOk
	}
	return
}
