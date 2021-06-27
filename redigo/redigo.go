package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/ratelimit"
)

func main() {
	flHost := flag.String("host", "", "redis host")
	flPort := flag.Int("port", 6379, "redis port")
	flAuth := flag.String("auth", "", "redis auth")
	flag.Parse()

	fmt.Printf("Start redigo test on %s:%d", *flHost, *flPort)

	rl := ratelimit.New(200000)

	pool := &redis.Pool{
		MaxIdle:   2000,
		MaxActive: 3000,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", *flHost, *flPort),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(500)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(500)*time.Millisecond))
		},
		Lifo: true,
		Wait: true,
	}
	defer pool.Close()
	for {
		rl.Take()
		go func() {
			conn := pool.Get()
			defer conn.Close()
			_, err := conn.Do("auth", *flAuth)
			if err != nil {
				fmt.Printf("err when auth! err: %v\n", err)
				os.Exit(0)
			}

			_, err = conn.Do("ping")
			if err != nil {
				fmt.Printf("err when ping! err: %v\n", err)
				os.Exit(0)
			}

			id, err := conn.Do("proxyid")
			if err != nil {
				fmt.Printf("err when ping! err: %v\n", err)
				os.Exit(0)
			}

			if string(id.([]byte)) == "187ecdcdc23cbd62c8970bbeb207443045f5bf77" {
				time.Sleep(time.Duration(rand.Int31n(200)+50) * time.Millisecond)
			} else {
				//time.Sleep(time.Duration(rand.Int31n(2)) * time.Millisecond)
			}
		}()
	}

	for {
		time.Sleep(time.Second)
	}
}
