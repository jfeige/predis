# predis
golang操作redis的实例，自定义了一个连接池


安装:

go get github.com/jfeige/predis




使用:

# 连接池配置
config := &PoolConfig{

	MaxCaps :100,
	MinCaps :10,
	IdleTimeout : 10*time.Second,
	Dial: func() (*Conn, error) {
		return Dial(network,address,pwd)
	},
	
}
# 初始化连接池
pool,err := NewPool(config)
if err != nil{

	fmt.Println(err)
	return
	
}
# 获取连接
conn,err := pool.Get()
if err != nil{
	//错误处理
}
defer pool.Put(conn)	//把连接重新放入池中

//do something

...

