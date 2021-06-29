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
		MaxIdle:   8000,
		MaxActive: 10000,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", *flHost, *flPort),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond))
		},
		SupportProxyId: true,
		Wait:           true,
	}
	defer pool.Close()
	for i := 0; i < 5000000; i++ {
		go func() {
			for {
				conn := pool.Get()
				_, err := conn.Do("auth", *flAuth)
				if err != nil {
					//	fmt.Printf("err when auth! err: %v\n", err)
					//	os.Exit(0)
				}

				_, err = conn.Do("ping")
				if err != nil {
					//	fmt.Printf("err when ping! err: %v\n", err)
					//	os.Exit(0)
				}

				_, err = conn.Do("hset", rand.Int31n(100), rand.Int31n(1000),
					fmt.Sprintf("%d:%d:%d:%d:%d:%d:%d:%d:%d:%d", rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31(), rand.Int31()))
				if err != nil {
					//	fmt.Printf("err when hset! err: %v\n", err)
					//os.Exit(0)
				}

				_, err = conn.Do("hgetall", rand.Int31n(100))
				if err != nil {
					//	fmt.Printf("err when hset! err: %v\n", err)
					//os.Exit(0)
				}
				conn.Close()
			}
		}()
	}

	for {
		time.Sleep(time.Second)
	}
}
