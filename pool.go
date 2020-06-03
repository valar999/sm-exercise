package pool

import (
	"log"
	"sync"
	"time"
)

type Connection interface {
	open()
	close()
}

type Conn struct {
	addr int32
}

func (c *Conn) open() {
	time.Sleep(time.Second)
	log.Println("open connection", c.addr)
}

func (c *Conn) close() {
	log.Println("close connection", c.addr)
}

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type connPool struct {
	sync.Mutex
	cache      map[int32]Connection
	cacheLocks map[int32]*sync.Mutex
}

func NewPool() Pool {
	return &connPool{
		cache:      make(map[int32]Connection),
		cacheLocks: make(map[int32]*sync.Mutex),
	}
}

func (pool *connPool) getConnection(addr int32) Connection {
	pool.Lock()
	conn, ok := pool.cache[addr]
	if ok {
		pool.Unlock()
	} else {
		mutex := new(sync.Mutex)
		pool.cacheLocks[addr] = mutex

		conn = &Conn{addr}
		mutex.Lock()
		pool.Unlock()
		conn.open()
		pool.cache[addr] = conn
		mutex.Unlock()
	}
	return conn
}

func (pool *connPool) onNewRemoteConnection(remotePeer int32, c Connection) {
	pool.Lock()
	mutex := pool.cacheLock[addr]
	mutex.Lock()
	if _, ok := pool.cache[remotePeer]; !ok {
		pool.cache[remotePeer] = c
	}
	mutex.Unlock()
	pool.Unlock()
}

func (pool *connPool) shutdown() {
	pool.Lock()
	defer pool.Unlock()
	for _, mutex := range pool.cacheLocks {
		mutex.Lock()
		defer mutex.Unlock()
	}
	for _, conn := range pool.cache {
		conn.close()
	}
}
