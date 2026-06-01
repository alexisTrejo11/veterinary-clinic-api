package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"clinic-vet-api/internal/shared/rediskeys"

	"github.com/redis/go-redis/v9"
)

type keyPrefixHook struct{}

func (keyPrefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (keyPrefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		prefixRedisCommand(cmd)
		return next(ctx, cmd)
	}
}

func (keyPrefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			prefixRedisCommand(cmd)
		}
		return next(ctx, cmds)
	}
}

func prefixRedisCommand(cmd redis.Cmder) {
	p := rediskeys.Prefix()
	if p == "" {
		return
	}

	args := cmd.Args()
	if len(args) < 2 {
		return
	}

	name := strings.ToUpper(fmt.Sprint(args[0]))
	switch name {
	case "GET", "SET", "SETEX", "SETNX", "GETEX", "GETDEL", "INCR", "INCRBY", "DECR", "DECRBY",
		"EXPIRE", "EXPIREAT", "TTL", "PTTL", "PERSIST", "TYPE", "RENAMENX", "EXISTS",
		"HGET", "HSET", "HMGET", "HMSET", "HGETALL", "HDEL", "HINCRBY", "HLEN", "HEXISTS":
		setArgKey(args, 1)
	case "DEL", "UNLINK", "TOUCH", "MGET":
		for i := 1; i < len(args); i++ {
			setArgKey(args, i)
		}
	case "RENAME":
		setArgKey(args, 1)
		setArgKey(args, 2)
	case "EVAL", "EVALSHA":
		nKeys := intArg(args[2])
		for i := 0; i < nKeys && 3+i < len(args); i++ {
			setArgKey(args, 3+i)
		}
	}
}

func setArgKey(args []interface{}, idx int) {
	if idx >= len(args) {
		return
	}
	key, ok := args[idx].(string)
	if !ok {
		return
	}
	args[idx] = rediskeys.Key(key)
}

func intArg(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	default:
		return 0
	}
}
