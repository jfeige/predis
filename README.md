# predis
golang操作redis的实例，自定义了一个连接池


安装:

go get github.com/jfeige/lredis




使用:

config := &PoolConfig{

	MaxCaps :100,
	MinCaps :10,
	IdleTimeout : 10*time.Second,
	Dial: func() (*Conn, error) {
		return Dial(network,address,pwd)
	},
}

pool,err := NewPool(config)
if err != nil{

	fmt.Println(err)
	return
}

conn,err := pool.Get()
