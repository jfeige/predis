package predis

import (
	"testing"
	"time"
	"fmt"
	"sync"
	"container/list"
)

var(
	network = "tcp"
	address = "127.0.0.1:6379"
	pwd = "lifei"
)

func Test_PoolList_One(t *testing.T){

	pool := &PoolList{
		Conns:list.New(),
		MinCaps:2,
		MaxCaps:10,
		IdleTimeout:10*time.Second,
		Dial: func() (*Conn, error) {
			return Dial(network,address,pwd)
		},
	}

	var wg sync.WaitGroup

	for i := 0;i < 30;i++{
		wg.Add(1)
		go testCmd(pool,i,&wg)
	}

	wg.Wait()

	fmt.Println("finish......")
}

func Test_PoolList_Two(t *testing.T){

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
	var wg sync.WaitGroup

	for i := 0;i < 30;i++{
		wg.Add(1)
		go testCmd(pool,i,&wg)
	}

	wg.Wait()

	fmt.Println("finish......")
}

func testCmd(pool *PoolList,position int,wg *sync.WaitGroup){

	conn,err := pool.GetConn()
	if err != nil{
		fmt.Println(err)
		return
	}
	defer func(){
		wg.Done()
		pool.Put(conn)
	}()

	tt := time.Now().UnixNano()

	dbsize,err := Int(conn.Cmd("dbsize"))

	fmt.Printf("协程:%d--------当前连接池中数量:%d---dbsize:%d,error:%v---%d\n",position,pool.GetCnt(),dbsize,err,tt)


}
