package predis

import (
	"sync"
	"time"
	"fmt"
	"bufio"
	"errors"
)

/**
	基于channel实现的连接池
 */


type PoolConfig struct {
	MaxCaps int
	MinCaps int
	IdleTimeout time.Duration
	Dial func()(*Conn,error)
}

type Pool struct {
	sync.Mutex
	conns chan *Conn
	maxCaps int
	minCaps int
	idleTimeout time.Duration
	dial func()(*Conn,error)
}


//初始化连接池
func NewPool(config *PoolConfig)(*Pool,error){
	pool := &Pool{
		conns:make(chan *Conn,config.MaxCaps),
		maxCaps:config.MaxCaps,
		minCaps:config.MinCaps,
		idleTimeout:config.IdleTimeout,
		dial:config.Dial,
	}
	for i := 0;i < config.MinCaps;i++{
		idleConn,err := pool.dial()
		if err != nil{
			return nil,fmt.Errorf("init redis pool has error:%",err)
		}

		//idleConn := &Conn{conn:conn, dst:make([]byte,0), br:bufio.NewReader(conn), bw:bufio.NewWriter(conn), pending:0,t:time.Now()}

		pool.conns <- idleConn
	}

	return pool,nil
}

//获取连接
func (this *Pool) Get()(*Conn,error){
	this.Lock()
	defer this.Unlock()

	if this.conns == nil{
		return nil,errors.New("connections has closed!")
	}
	for{
		select{
		case conn :=<- this.conns:
			if timeout := this.idleTimeout;timeout > 0{
				if conn.t.Add(timeout).Before(time.Now()){
					//该连接已超时，关闭
					this.close(conn)
					continue
				}
			}
			return conn,nil
		default:
			idleConn,err := this.dial()
			if err != nil{
				return nil,err
			}
			//idleConn := &Conn{conn:conn, dst:make([]byte,0), br:bufio.NewReader(conn), bw:bufio.NewWriter(conn), pending:0,t:time.Now()}

			return idleConn,nil
		}
	}
}

//把连接重新放入连接池
func (this *Pool) Put(c *Conn)error{
	this.Lock()
	defer this.Unlock()

	if this.conns == nil{
		return fmt.Errorf("connection pool is nil!")
	}

	if len(this.conns) >= this.maxCaps{
		return nil
	}
	idleConn := &Conn{conn:c.conn, dst:make([]byte,0), br:bufio.NewReader(c.conn), bw:bufio.NewWriter(c.conn), pending:0,t:time.Now()}

	this.conns <- idleConn

	return nil
}

//关闭连接
func (this *Pool) close(conn *Conn){
	conn.conn.Close()
}


//当前连接池中的连接数
func (this *Pool) Conns()int{
	this.Lock()
	defer this.Unlock()

	if this.conns == nil{
		return 0
	}
	return len(this.conns)
}

