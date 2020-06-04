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
		if c.conn == nil {
			log.Println("get cache, but no conn")
			c.Lock()
			c.Unlock()
		}
	} else {
		c = &Conn{
			ch:   make(chan Connection, 2),
		}
		pool.cache[addr] = c

		c.Lock()
		pool.Unlock()

		openConn := NewConn(addr, openDelay)
		log.Println(openConn.n, "conn")
		go func(ch chan Connection) {
			log.Println(openConn.n, "open1")
			openConn.open()
			log.Println(openConn.n, "open2")
			ch <- openConn
		}(c.ch)
		select {
		case conn := <-c.ch:
			pool.Lock()
			pool.cache[addr].conn = conn
			pool.Unlock()
		}
		c.Unlock()
	}
	log.Println(c.conn.(*connMock).n, "ret")
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
