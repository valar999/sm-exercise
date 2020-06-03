package pool

type Connection interface {
	open()
	close()
}

type Conn struct {
	addr int32
}

func (c *Conn) open() {
	// time.Sleep(time.Second)
	// log.Println("open connection", c.addr)
}

func (c *Conn) close() {
	// log.Println("close connection", c.addr)
}
