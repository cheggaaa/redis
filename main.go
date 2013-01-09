package redis

func NewClient(addr string) (client *Redis, err error) {
	config := &Config{}
	config.Addr = addr
	return NewClientFromConfig(config)
}

func NewClientFromConfig(config *Config) (client *Redis, err error) {
	err = config.Init()
	if err != nil {
		return
	}
	client = &Redis{
		Config : config,
	}
	err = client.connect()
	if err != nil {
		client = nil
	}
	return
}