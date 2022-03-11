package redis_lock

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestRedisLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	rl, err := NewRedisLock(client, "mykey")
	rl.SetExpire(3)
	defer func() {
		ok, err := rl.Release()
		if err != nil {
			fmt.Println("release lock err: ", err)
		}
		if !ok {
			fmt.Println("release lock failed")
		} else {
			fmt.Println("release lock success")
		}
	}()
	if err != nil {
		panic(err)
	}
	ok, err := rl.Acquire()
	if err != nil {
		panic(err)
	}
	if ok {
		fmt.Println("acquire lock success")
		// dothings
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("acquire lock failed")
	}
}
