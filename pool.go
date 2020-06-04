package pool

import (
	"log"
	"sync"
	"time"
)

var openDelay = time.Millisecond * 100

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type Conn struct {
	sync.Mutex
	ch   chan Connection
	conn Connection
}

type pool struct {
	sync.Mutex
	cache      map[int32]*Conn
	isShutdown bool
}

func NewPool() Pool {
	return &pool{
		cache: make(map[int32]*Conn),
	}
}

func (pool *pool) getConnection(addr int32) Connection {
	pool.Lock()
	if pool.isShutdown {
		pool.Unlock()
		return nil
	}
	c, ok := pool.cache[addr]
	if ok {
		pool.Unlock()
	} else {
		c = &Conn{
			conn: NewConn(addr, openDelay),
			ch:   make(chan Connection, 2),
		}
		pool.cache[addr] = c

		c.Lock()
		pool.Unlock()
		log.Println(c.conn.(*connMock).n, "conn")
		go func(ch chan Connection) {
			log.Println(c.conn.(*connMock).n, "open1")
			c.conn.open()
			log.Println(c.conn.(*connMock).n, "open2")
			ch <- c.conn
		}(c.ch)
		select {
		case conn := <-c.ch:
			log.Println(conn.(*connMock).n, "ret")
			c.conn = conn
		}
		c.Unlock()
	}
	return c.conn
}

func (pool *pool) onNewRemoteConnection(remotePeer int32, c Connection) {
	pool.Lock()
	defer pool.Unlock()
	if pool.isShutdown {
		return
	}
	conn, ok := pool.cache[remotePeer]
	if ok {
		log.Println(c.(*connMock).n, "new ch<-")
		select {
		case conn.ch <- c:
		default:
		}
		log.Println(c.(*connMock).n, "new2 ch<-")
	} else {
		log.Println(c.(*connMock).n, "new store")
		pool.cache[remotePeer] = &Conn{conn: c}
	}
}

func (pool *pool) shutdown() {
	pool.Lock()
	defer pool.Unlock()
	for _, c := range pool.cache {
		c.Lock()
		defer c.Unlock()
	}
	for _, c := range pool.cache {
		c.conn.close()
	}
	pool.isShutdown = true
}
