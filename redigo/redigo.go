package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	flHost := flag.String("host", "", "redis host")
	flPort := flag.Int("port", 6379, "redis port")
	flAuth := flag.String("auth", "", "redis auth")
	flag.Parse()

	fmt.Printf("Start redigo test on %s:%d", *flHost, *flPort)

	pool := &redis.Pool{
		MaxIdle:   4000,
		MaxActive: 5000,
		Dial: func() (redis.Conn, error) {
			c, connErr := redis.Dial("tcp", fmt.Sprintf("%s:%d", *flHost, *flPort),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(500)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(500)*time.Millisecond))
			if connErr != nil {
				return nil, connErr
			}

			_, authErr := c.Do("auth", *flAuth)
			if authErr != nil {
				c.Close()
				return nil, authErr
			}
			return c, nil
		},
		Lifo: false,
		Wait: true,
	}
	defer pool.Close()
	for i := 0; i < 1000000; i++ {
		go func() {
			for {
				conn := pool.Get()
				_, _ = conn.Do("hset", rand.Int31n(100), rand.Int31n(1000), rand.Int31())
				time.Sleep(time.Duration(rand.Int63n(5)) * time.Millisecond)
				conn.Close()
			}
		}()
	}

	for {
		time.Sleep(time.Second)
	}
}
