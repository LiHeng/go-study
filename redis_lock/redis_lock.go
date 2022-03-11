package redis_lock

import (
	"strconv"
	"sync/atomic"

	"github.com/go-redis/redis"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
)

// 可重入
const (
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	unlockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

const (
	millisPerSecond = 1000
	tolerance       = 500 // milliseconds
)

type RedisLock struct {
	store *redis.Client
	// 超时时间
	seconds uint32
	// 锁key
	key string
	// 锁value，防止锁被别人获取到
	id string
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *redis.Client, key string) (*RedisLock, error) {
	uid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &RedisLock{
		store: store,
		key:   key,
		// 获取锁时，锁的值通过随机字符串生成
		id: uid,
	}, nil
}

func (rl *RedisLock) Acquire() (bool, error) {
	// 获取过期时间
	seconds := atomic.LoadUint32(&rl.seconds)
	// 默认锁过期时间为500ms，防止死锁
	resp, err := rl.store.Eval(lockCommand, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	}).Result()

	// key不存在也代表失败
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		logrus.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	logrus.Errorf("Unknown reply when acquiring lock for %s: %v", rl.key, resp)
	return false, nil
}

func (rl *RedisLock) Release() (bool, error) {
	resp, err := rl.store.Eval(unlockCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}

// SetExpire sets the expire.
// 需要注意的是需要在Acquire()之前调用
// 不然默认为500ms自动释放
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}
