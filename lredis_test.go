package predis

import (
	"testing"
	"fmt"
	"time"
)

func Test_Send(t *testing.T) {
	var(
		network = "tcp"
		address = "127.0.0.1:6379"
		pwd = "lifei"
	)
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
	fmt.Printf("--1--当前连接池中数量:%d\n",pool.Conns())
	conn,err := pool.Get()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("--2--当前连接池中数量:%d\n",pool.Conns())
	defer func(){
		pool.Put(conn)	//很重要，把连接重新放入连接池
		fmt.Printf("--3--当前连接池中数量:%d\n",pool.Conns())
	}()			
	dbsize,err := Int(conn.Cmd("dbsize"))

	fmt.Println(dbsize,err)



	time.Sleep(3*time.Second)

}



