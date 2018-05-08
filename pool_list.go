package predis

import (
	"sync"
	"container/list"
	"time"
	"fmt"
	"bufio"
)

/**
	基于list实现的连接池
 */

 type PoolList_Config struct {
 	MinCaps int
 	MaxCaps int
 	IdleTimeout time.Duration
 	Created bool				//是否在初始化连接池时，生成最小连接数
 	Dial func()(*Conn,error)
 }

 type PoolList struct {
 	sync.Mutex
 	Conns *list.List
 	MinCaps int
 	MaxCaps int
 	IdleTimeout time.Duration
 	Dial func()(*Conn,error)
 }


 //初始化连接池
func NewPoolList(config *PoolList_Config)(*PoolList,error){
	pool := &PoolList{
		MinCaps:config.MinCaps,
		MaxCaps:config.MaxCaps,
		IdleTimeout:config.IdleTimeout,
		Conns:list.New(),
		Dial:config.Dial,
	}
	if pool.MaxCaps > 0 && pool.MinCaps > pool.MaxCaps{
		return nil,fmt.Errorf("predis:configuration parameter error.mincaps:%d,maxcaps:%d",pool.MinCaps,pool.MaxCaps)
	}
	if config.Created{
		for i := 0; i < config.MinCaps;i++{
			c,err := pool.Dial()
			if err != nil{
				pool.Release()
				return nil,err
			}
			pool.Conns.PushFront(c)
		}
	}

	return pool,nil
}


 //获取连接
 func (this *PoolList)GetConn()(*Conn,error){
 	this.Lock()
 	defer this.Unlock()

 	for e := this.Conns.Back();e !=nil;e.Next(){
		c := e.Value.(*Conn)
 		if timeout := this.IdleTimeout;timeout > 0{
			if c.t.Add(this.IdleTimeout).Before(time.Now()){
				//连接已超时,关闭该连接
				this.Conns.Remove(e)
				this.close(c)
				continue
			}
		}
		this.Conns.Remove(e)
		return c,nil
	}


	//生成新的连接
	c,err := this.Dial()
	if err != nil{
		return nil,err
	}
	return c,nil
 }



 //把连接放回池中
func (this *PoolList) Put(c *Conn){
	this.Lock()
	defer this.Unlock()

	if this.MaxCaps > 0 &&  this.Conns.Len() >= this.MaxCaps{
		return
	}
	c.bw = bufio.NewWriter(c.conn)
	c.br = bufio.NewReader(c.conn)
	c.t = time.Now()
	c.pending = 0
	this.Conns.PushFront(c)
}



 //得到当前连接池中可用连接数量
 func (this *PoolList) GetCnt()int{
 	this.Lock()
 	defer this.Unlock()
 	return this.Conns.Len()
 }


 //关闭连接
 func (this *PoolList) close(c *Conn){
	c.conn.Close()
 }


 //释放所有连接
 func (this *PoolList) Release(){
 	this.Lock()
 	defer this.Lock()

 	if this.Conns == nil{
 		return
	}
	for e := this.Conns.Front();e != nil;e.Next(){
		idleConn := e.Value.(*Conn)
		idleConn.conn.Close()
		this.Conns.Remove(e)
	}
	return
 }