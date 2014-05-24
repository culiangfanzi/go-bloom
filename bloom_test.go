package bloom

import (
	"fmt"
	"github.com/kristinn/redigo/redis"
	"testing"
)

func newRedisConnection() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func TestRedisInit(t *testing.T) {
	pool := redis.NewPool(newRedisConnection, 5)
	defer pool.Close()

	_, err := NewRedis(pool, "redis-init-test", 15000, 7)
	if err != nil {
		t.Fatal(err)
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Do("FLUSHALL")
}

func TestRedisSave(t *testing.T) {
	pool := redis.NewPool(newRedisConnection, 5)
	defer pool.Close()

	r, err := NewRedis(pool, "redis-save-test", 15000, 7)
	if err != nil {
		t.Fatal(err)
	}

	r.Append([]byte("afi"))
	r.Save()

	exists, err := r.Exists([]byte("afi"))
	if !exists {
		t.Fatal("afi should exist in the Redis backend")
	}
	if err != nil {
		t.Fatal(err)
	}

	exists, err = r.Exists([]byte("amma"))
	if exists {
		t.Fatal("amma shouldn't exist in the Redis backend")
	}
	if err != nil {
		t.Fatal(err)
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Do("FLUSHALL")
}

func TestBitsetSave(t *testing.T) {
	b := NewBitset(15000, 7)

	b.Append([]byte("afi"))
	b.Save()

	exists, err := b.Exists([]byte("afi"))
	if !exists {
		t.Fatal("afi should exist in the Bitset backend")
	}
	if err != nil {
		t.Fatal(err)
	}

	exists, err = b.Exists([]byte("amma"))
	if exists {
		t.Fatal("amma shouldn't exist in the Bitset backend")
	}
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkRedisQueueAppend(b *testing.B) {
	pool := redis.NewPool(newRedisConnection, 7)
	defer pool.Close()

	r, err := NewRedis(pool, "redis-queue-append-benchmark", 15000, 7)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		r.Append([]byte(fmt.Sprintf("afi.%d", i)))
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Do("FLUSHALL")
}

func BenchmarkRedisSave(b *testing.B) {
	pool := redis.NewPool(newRedisConnection, 7)
	defer pool.Close()

	r, err := NewRedis(pool, "redis-save-benchmark", 15000, 7)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		r.Append([]byte(fmt.Sprintf("afi.%d", i)))
		r.Save()
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Do("FLUSHALL")
}

func BenchmarkRedisExists(b *testing.B) {
	pool := redis.NewPool(newRedisConnection, 7)
	defer pool.Close()

	r, err := NewRedis(pool, "redis-exists-benchmark", 15000, 7)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 15000; i++ {
		r.Append([]byte(fmt.Sprintf("afi.%d", i)))
	}
	r.Save()

	for i := 0; i < b.N; i++ {
		r.Exists([]byte("afi.7500"))
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Do("FLUSHALL")
}

func BenchmarkBitsetAppend(b *testing.B) {
	bits := NewBitset(15000, 7)

	for i := 0; i < b.N; i++ {
		bits.Append([]byte(fmt.Sprintf("afi.%d", i)))
	}
}

func BenchmarkBitsetSave(b *testing.B) {
	bits := NewBitset(15000, 7)

	for i := 0; i < b.N; i++ {
		bits.Append([]byte(fmt.Sprintf("afi.%d", i)))
		bits.Save()
	}
}

func BenchmarkBitsetExists(b *testing.B) {
	bits := NewBitset(15000, 7)

	for i := 0; i < 15000; i++ {
		bits.Append([]byte(fmt.Sprintf("afi.%d", i)))
	}
	bits.Save()

	for i := 0; i < b.N; i++ {
		bits.Exists([]byte("afi.7500"))
	}
}
