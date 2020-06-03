package main

import (
	"log"
	"sync"
	"time"
)

type Connection interface {
	open()
	close()
}

type Conn struct{
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
	shutdown()
}

type connPool struct {
	sync.Map
}

func NewPool() Pool {
	return &connPool{}
}

func (pool *connPool) getConnection(addr int32) Connection {
	conn, ok := pool.Load(addr)
	if ok {
		return conn.(Connection)
	} else {
		conn := &Conn{addr}
		conn.open()
		pool.Store(addr, conn)
		return conn
	}
}

func (pool *connPool) shutdown() {
}

func main() {
	pool := NewPool()
	log.Println(pool.getConnection(1))
	log.Println(pool.getConnection(2))
	log.Println(pool.getConnection(1))
}
