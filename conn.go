package pool

import "time"

type Connection interface {
	open()
	close()
}

type conn struct {
	addr int32
}

func (c *conn) open() {
	time.Sleep(time.Microsecond * 100)
	// log.Println("open connection", c.addr)
}

func (c *conn) close() {
	// log.Println("close connection", c.addr)
}
