package pool

import (
	"sync"
)

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type connPool struct {
	sync.Mutex
	isShutdown bool
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
		pool.Lock()
		pool.cache[addr] = conn
		pool.Unlock()
		mutex.Unlock()
	}
	return conn
}

func (pool *connPool) onNewRemoteConnection(remotePeer int32, c Connection) {
	pool.Lock()
	mutex := pool.cacheLocks[remotePeer]
	mutex.Lock()
	if _, ok := pool.cache[remotePeer]; !ok {
		pool.cache[remotePeer] = c
	}
	mutex.Unlock()
	pool.Unlock()
}

func (pool *connPool) shutdown() {
	pool.Lock()
	pool.isShutdown = true
	defer pool.Unlock()
	for _, mutex := range pool.cacheLocks {
		mutex.Lock()
		defer mutex.Unlock()
	}
	for _, conn := range pool.cache {
		conn.close()
	}
}
