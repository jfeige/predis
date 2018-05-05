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
 	conns *list.List
 	mincaps int
 	maxcaps int
 	idletimeout time.Duration
 	dial func()(*Conn,error)
 }


 //初始化连接池
func NewPoolList(config *PoolList_Config)(*PoolList,error){
	pool := &PoolList{
		mincaps:config.MinCaps,
		maxcaps:config.MaxCaps,
		idletimeout:config.IdleTimeout,
		conns:list.New(),
		dial:config.Dial,
	}
	if pool.maxcaps > 0 && pool.mincaps > pool.maxcaps{
		return nil,fmt.Errorf("predis:configuration parameter error.mincaps:%d,maxcaps:%d",pool.mincaps,pool.maxcaps)
	}
	if config.Created{
		for i := 0; i < config.MinCaps;i++{
			c,err := pool.dial()
			if err != nil{
				pool.Release()
				return nil,err
			}
			pool.conns.PushFront(c)
		}
	}

	return pool,nil
}


 //获取连接
 func (this *PoolList)GetConn()(*Conn,error){
 	this.Lock()
 	defer this.Unlock()

 	for e := this.conns.Back();e !=nil;e.Next(){
		c := e.Value.(*Conn)
 		if timeout := this.idletimeout;timeout > 0{
			if c.t.Add(this.idletimeout).Before(time.Now()){
				//连接已超时,关闭该连接
				this.conns.Remove(e)
				this.close(c)
				continue
			}
		}
		this.conns.Remove(e)
		return c,nil
	}


	//生成新的连接
	c,err := this.dial()
	if err != nil{
		return nil,err
	}
	return c,nil
 }



 //把连接放回池中
func (this *PoolList) Put(c *Conn){
	this.Lock()
	defer this.Unlock()

	if this.maxcaps > 0 &&  this.conns.Len() >= this.maxcaps{
		return
	}
	c.bw = bufio.NewWriter(c.conn)
	c.br = bufio.NewReader(c.conn)
	c.t = time.Now()
	c.pending = 0
	this.conns.PushFront(c)
}



 //得到当前连接池中可用连接数量
 func (this *PoolList) GetCnt()int{
 	this.Lock()
 	defer this.Unlock()
 	return this.conns.Len()
 }


 //关闭连接
 func (this *PoolList) close(c *Conn){
	c.conn.Close()
 }


 //释放所有连接
 func (this *PoolList) Release(){
 	this.Lock()
 	defer this.Lock()

 	if this.conns == nil{
 		return
	}
	for e := this.conns.Front();e != nil;e.Next(){
		idleConn := e.Value.(*Conn)
		idleConn.conn.Close()
		this.conns.Remove(e)
	}
	return
 }