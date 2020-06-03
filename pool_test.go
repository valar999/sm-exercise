package pool

import (
	"log"
	"testing"
)

func TestConnection(t *testing.T) {
	pool := NewPool()
	log.Println(pool.getConnection(1))
	log.Println(pool.getConnection(2))
	log.Println(pool.getConnection(1))
	pool.shutdown()
}
