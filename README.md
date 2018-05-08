# predis
golang操作redis的实例，实现了一个连接池<br/>其中socket通讯和解析这块，是直接拿之前写的lredis代码，凑合着看吧，后面会继续完善<br>连接池这块写了两个，一个通过List实现，一个通过channel来实现。具体的实现，有略微不同，channel是在初始化时，就往池中放入一定的连接，而List则是在获取连接时，动态创建连接，使用完毕后，再放入连接池

## 安装:

```
go get github.com/jfeige/predis
```


## 使用:


### 初始化连接池
```
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
```
##### or
```
config := &PoolList_Config{
	MinCaps:2,
	MaxCaps:10,
	IdleTimeout:10*time.Second,
	Dial: func() (*Conn, error) {
		return Dial(network,address,pwd)
	},
}
pool,err := NewPoolList(config)
if err != nil{
	fmt.Println(err)
	return
}
```
##### or
```
pool := &PoolList{
	Conns:list.New(),
	MinCaps:2,
	MaxCaps:10,
	IdleTimeout:10*time.Second,
	Dial: func() (*Conn, error) {
		return Dial(network,address,pwd)
	},
}
```

### 获取连接
```
conn,err := pool.Get()
if err != nil{
	//错误处理
}
defer pool.Put(conn)	//把连接重新放入池中
...
//do something
```

##### or

```
conn,err := pool.GetConn()
if err != nil{
	fmt.Println(err)
	return
}
defer func(){
	pool.Put(conn) //把连接重新放入池中
}()
...
//do something
```
