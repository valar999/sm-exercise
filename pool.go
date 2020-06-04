package pool

import (
	"sync"
	"log"
	"time"
)

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type Conn struct {
	sync.Mutex
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
		c = &Conn{conn: &conn{addr, time.Microsecond * 1000}}
		pool.cache[addr] = c

		c.Lock()
		pool.Unlock()
		log.Println("open")
		c.conn.open()
		log.Println("open2")
		c.Unlock()
	}
	return c.conn
}

func (pool *pool) onNewRemoteConnection(remotePeer int32, c Connection) {
	log.Println("onNewRemoteConnection")
	pool.Lock()
	defer pool.Unlock()
	if pool.isShutdown {
		return
	}
	_, ok := pool.cache[remotePeer]
	if !ok {
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
